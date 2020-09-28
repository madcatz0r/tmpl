# tmpl
##postgresql templator

###Шаблонизатор sql запросов pre-alfa

Порядок использования:

1. Покрытие тэгами моделей
    ```go
    type ExampleStruct struct {
   	    ID              uuid.UUID `tmpl:"type=primary"`
   	    OuterID         int64     `tmpl:"type=outer"`
   	    Amount          int64     `tmpl:"type=upsert"`
   	    AnotherAmount   int64     `tmpl:"type=upsert"`
   	    SomeUnusedThing int64     `tmpl:"-"`
   	    SomeUsefulThing string    `tmpl:"type=upsert"`
   	    CreatedAt       time.Time `tmpl:"type=insert,default=now() AT TIME ZONE 'UTC'"`
   	    UpdatedAt       time.Time `tmpl:"type=must_upd,default=now() AT TIME ZONE 'UTC'"`
    } 
    ```
1. (необязательно) c помощью tmplgen -model=<models path> -out=<vars path> сгенерировать ассоциации свойств моделей на поля БД
model CamelCase -> db snake_case
    ```go
    package ExampleStruct
   
    const ( 
        ID              = "example_struct.id"
        OuterID         = "example_struct.outer_id"
        Amount          = "example_struct.amount"
        AnotherAmount   = "example_struct.another_amount"
        SomeUnusedThing = "example_struct.some_unused_thing"
        SomeUsefulThing = "example_struct.some_useful_thing"
        CreatedAt       = "example_struct.created_at"
        UpdatedAt       = "example_struct.updated_at"
    )
    ```
1. На старте обязательно выполнить err := ParseTags{ExampleStruct{})
1. Применение шаблонов
   ```go
   example := &ExampleStruct{}
   
   insertQuery, orderedValues := tmpl.Insert(example)
   updateQuery, orderedValues := tmpl.Update(example)
   partialUpdateQuery, orderedValues := tmpl.Update(example, ExampleStruct.Amount)
   upsertQuery, orderedValues := tmpl.Upsert(example)
   partialUpdateUpsertQuery, orderedValues := tmpl.Upsert(example, ExampleStruct.Amount)
   
   query := Select(ExampleStruct.Amount).From(ExampleStruct{}).Where(conditions.Eq(ExampleStruct.ID, conditions.String("1")))
   // больше примеров использования в select_test.go
   ```

TODO:
- [ ] добавить модуль example
- [ ] переработать и перевести README
