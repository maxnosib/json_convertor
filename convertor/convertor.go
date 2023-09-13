package convertor

import (
	"log"
	"reflect"
)

/*
UnmarshalMap функция которая обрабатявает любой входящий json в нужную структуру
входящий json расскладывается в одномерный массив через ф-ю DimensionalMap
затем проходя по заданной структуре читаются теги ref и занчения ищятся в входящем массиве
затем если данные найдены то записываются в структуру
в структуре не должно быть типа map[string] т.к. этот тип будет разложен на более простые в DimensionalMap
функция UnmarshalMap преобразует базовые типы (string int float bool), а так же преобразует к указателю и слайсу (с простыми типами)
поля структуры должны быть экспортируемыми и иметь тег ref
*/

func UnmarshalMap(structure any, data map[string]interface{}) error {
	refVal := reflect.ValueOf(structure).Elem() // получаем значение
	refType := refVal.Type()                    // получаем тип

	for i := 0; i < refVal.NumField(); i++ { // итерируемся по полям структуры
		// получаем поле структуры
		refField := refVal.Field(i)

		// проверяем экспортируемое ли поле или нет чтоб не словить панику
		if !refType.Field(i).IsExported() {
			log.Printf("field %s not exported", refType.Field(i).Name)
			continue
		}

		// проверяем что поле валидно
		if !refField.IsValid() {
			log.Printf("no valud field: %s", refType.Field(i).Name)
			continue
		}

		// проверяем что в поле можно записывать
		if !refField.CanSet() {
			log.Printf("cannot set value in field: %s", refType.Field(i).Name)
			continue
		}

		// проверяем если структура то проходимся  по ней
		if refField.Kind() == reflect.Struct {
			str := refField.Addr().Interface() // получаем адрес поля чтоб передать его рекурсию
			UnmarshalMap(str, data)
			continue
		}

		// проверки закончены начинаем работать уже с полями

		// получаем тег
		tag := refType.Field(i).Tag.Get("ref")

		// получаем значение из входящего jsona
		value := data[tag]
		if value == nil {
			log.Printf("no data from ref tag: %s  from field: %s", tag, refType.Field(i).Name)
			continue
		}

		// получаем значение из под интерфейса
		val := reflect.ValueOf(value)

		// получаем тип поля
		if refField.Kind() == reflect.Ptr {
			// т.к. поле указатель то создаем указатель на значение
			ptrVal := reflect.New(val.Type())

			// делаем инициализацию поля
			refField.Set(ptrVal)

			// получаем значение на которое указывает refField мы его уже проинициализировали чтоб было не nil
			refField = reflect.Indirect(refField)

			// записываем значение в поле структуры
			refField.Set(reflect.ValueOf(value))
			continue
		}

		// получаем тип поля
		if refField.Kind() == reflect.Slice {
			// создаем пустой слайс чтоб записать его после в поле структуры
			slice := reflect.MakeSlice(refField.Type(), 0, val.Len())

			// проходимся по входящему слайсу
			for i := 0; i < val.Len(); i++ {
				// val.Index(i).Interface() превращяет reflect.Value в interface
				// через reflect.ValueOf создаем новое reflect.Value которое можно будет конвертировать в нужный тип
				elem := reflect.ValueOf(val.Index(i).Interface())

				// конвертируем элемент массива в тип данных (получаем через refField.Type().Elem()) лежащий в слайсе
				// затем добавляем этот элемент в слайс
				slice = reflect.Append(slice, elem.Convert(refField.Type().Elem()))

			}

			// записываем значение полученный слайс в поле структуры
			refField.Set(slice)
			continue
		}

		// проверяем можем ли мы конвертировать тип данных в тип поля структуры
		if !val.CanConvert(refField.Type()) {
			log.Printf("can not convert field: %s", refType.Field(i).Name)
			continue
		}

		// записываем значение в структуру  конвертируя его в тип поля структуры
		refField.Set(val.Convert(refField.Type()))

	}

	return nil
}

// делам мапу одномерной
func DimensionalMap(in, out map[string]interface{}) {
	for key, val := range in {
		if new, ok := val.(map[string]interface{}); ok {
			DimensionalMap(new, out)
			continue
		}

		out[key] = val
	}
}
