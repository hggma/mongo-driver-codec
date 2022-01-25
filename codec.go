package mongo_driver_codec

import (
	"fmt"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

// Time
// Usage:
// 		mongo.Connect(context.Background(),
//					options.Client().ApplyURI("mongodb://127.0.0.1:27017").
//					SetRegistry(bson.NewRegistryBuilder().
//					RegisterTypeDecoder(reflect.TypeOf(time.Time{}), Time{}).
//					Build()))
type Time time.Time

func (t Time) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	timeType := reflect.TypeOf(time.Time{})
	if !val.IsValid() || !val.CanSet() || val.Type() != timeType {
		return bsoncodec.ValueDecoderError{
			Name:     "timeDecodeValue",
			Types:    []reflect.Type{timeType},
			Received: val,
		}
	}
	if vr.Type() != bsontype.DateTime {
		return fmt.Errorf("received invalid BSON type to decode into time.Time: %s", vr.Type())
	}
	dec, err := vr.ReadDateTime()
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(time.Unix(dec/1e3, dec%1e3*1e6)))
	return nil
}

// Decimal
// Usage:
// 		mongo.Connect(context.Background(),
//					options.Client().ApplyURI("mongodb://127.0.0.1:27017").SetRegistry(bson.NewRegistryBuilder().
//					RegisterTypeDecoder(reflect.TypeOf(decimal.Decimal{}), Decimal{}).
//					RegisterTypeEncoder(reflect.TypeOf(decimal.Decimal{}), Decimal{}).
//					Build()))
type Decimal decimal.Decimal

func (d Decimal) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	decimalType := reflect.TypeOf(decimal.Decimal{})
	if !val.IsValid() || !val.CanSet() || val.Type() != decimalType {
		return bsoncodec.ValueDecoderError{
			Name:     "decimalDecodeValue",
			Types:    []reflect.Type{decimalType},
			Received: val,
		}
	}
	var value decimal.Decimal
	switch vr.Type() {
	case bsontype.Decimal128:
		dec, err := vr.ReadDecimal128()
		if err != nil {
			return err
		}
		value, err = decimal.NewFromString(dec.String())
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("received invalid BSON type to decode into decimal.Decimal: %s", vr.Type())
	}
	val.Set(reflect.ValueOf(value))
	return nil
}

func (d Decimal) EncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	decimalType := reflect.TypeOf(decimal.Decimal{})
	if !val.IsValid() || val.Type() != decimalType {
		return bsoncodec.ValueEncoderError{
			Name:     "decimalEncodeValue",
			Types:    []reflect.Type{decimalType},
			Received: val,
		}
	}
	dec := val.Interface().(decimal.Decimal)
	dec128, err := primitive.ParseDecimal128(dec.String())
	if err != nil {
		return err
	}
	return vw.WriteDecimal128(dec128)
}
