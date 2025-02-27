# Расширение message и enum типов для работы с прото
## Запуск
- Запустить можно тупо `go run main.go` из корня, что бы посмотреть вывод

## Описание
- Этот репозиторий отражает работу с расширениями для прото сообщений и энум типов

## Что происходит
Функция main запускает логику, ничего особенного

### GetEnumMessage и GetMessage
- Это функции, которые формируют mock сообщение в []byte . Имитируют безликое сообщение, которое получил клиент
- GetEnumMessage формирует сообщение, внутри которого лежит enum тип
- GetMessage формирует обычное сообщение типа Message без уникальных полей
- GetMessage возвращает сообщение, которое содержит сообщение, и его proto дескриптор

### ServeMessage и ServeEnumMessage
- функции которые обрабатывают полученное сообщение в []byte
- ServeMessage тащит и записывает в stdout инфу о типе пришедшего сообщения из дескриптора
- ServeEnumMessage делает то же самое, только для сообщения с enum полем

### proto/test.proto
#### google.protobuf.DescriptorProto
- К сообщению, помимо полей этого сообщения мы можем добавить поле типа `google.protobuf.DescriptorProto`
- Для того , что бы этот тип добавить, нужно импортировать `google/protobuf/descriptor.proto`
- С таким подходом, на серверной стороне (там, где формируется proto сообщение), необходимо явно добавить к сообщению нужный дескриптор.
- Это не всегда удобно, особенно, если у нас нет доступа к серверной части (но мы можем изменить контракт)
#### extensions
- Расширения типов - хороший инструмент для обогащения любых типов какими-либо метаданными.
- Например, мы можем добавить описание к enum значению, с помощью расширения для `FieldOptions`
```
extend google.protobuf.FieldOptions {
  string my_custom_option = 50000;
}

message Person {
  string name = 1 [(my_custom_option) = "some metadata"];
  int32 age = 2;
}
```
- нужно указать тип опции, и уникальный номер поля (можно просто побольше указать, у меня тут 50000 например)
- Стоит учитывать, что такой экстеншн неприменим к enum значениям. Что бы мы могли расширить enum значение, нужно указать другой тип экстеншена:
```
extend google.protobuf.EnumValueOptions {
  optional string abbr = 54321;
  optional google.protobuf.DescriptorProto descriptor = 54322;
}
```
- в этом случае я добавил экстеншен для EnumValue с двумя полями - abbr, аббревиатуру значения и дескриптор relied структуры
- Про дескриптор relied структуры - теоретическая ситуация:
  - Мы - клиент. Получаем сообщение, в котором enum значение
  - Мы хотим в зависимости от типа этого сигнала, сформировать новое сообщение, и транслировать его куда то дальше, по нашему флоу
  - так как у нас enum значение расширено дескриптором - мы можем вытащить этот дескриптор с помощью GetExtension:
```go
    messageDescriptor := proto.GetExtension(options, test.E_Descriptor).(*descriptorpb.DescriptorProto)
	if messageDescriptor == nil {
		fmt.Println("Message descriptor not found for enum value:", enumValue.Name())
		return nil
	}
```
- Для того , что бы снять ответственность с серверной части - нам стоит явно указать дескриптор и abbr для энум значения:
```
message PersonsEnumType {
  enum SignalType {
    UNDEFINED = 0 [
      (abbr)="undefined_topic",
      (descriptor) = {
        name: "test.Undefined",
        field: [
          { name: "info", number: 1, type: TYPE_STRING }
        ]
      }
    ];
  }
}

message Undefined {
  string info = 1;
}
```
- Делается это таким образом. У нас в прото файле есть какая то структура UNDEFINED - с которой мы хотим ассоциировать этот энум
- Мы явно указываем в поле `(descriptor)` структуру ассоциированного типа, таким образом, что бы protoc смог это прочитать и сгенерировать эксеншн (что бы с ним можно было в коде работать)


