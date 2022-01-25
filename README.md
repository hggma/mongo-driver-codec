# mongo-driver-codec
mongodb 类型转换，目前支持decimal和时区差。

## Quick Start

```go
// 在项目中引入 codec
import "github.com/hggma/mongo-driver-codec"

// 注入codec既可
client, err := mongo.Connect(context.Background(),
    options.Client().ApplyURI("mongodb://127.0.0.1:27017").SetConnectTimeout(2*time.Second).
    SetRegistry(bson.NewRegistryBuilder().
    RegisterTypeDecoder(reflect.TypeOf(time.Time{}), mongo_driver_codec.Time{}).
    RegisterTypeDecoder(reflect.TypeOf(decimal.Decimal{}), mongo_driver_codec.Decimal{}).
    RegisterTypeEncoder(reflect.TypeOf(decimal.Decimal{}), mongo_driver_codec.Decimal{}).
    Build()))

```