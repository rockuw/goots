goots
=====

Aliyun OTS(Open Table Service) golang SDK.

## Support API
- **Table**
	- [CreateTable](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/CreateTable.md)
	- DeleteTable
	- [ListTable](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/ListTable.md)
	- UpdateTable
	- DescribeTable
- **SingleRow**
	- GetRow
	- PutRow
	- UpdateRow
	- DeleteRow
- **BatchRow**
	- BatchGetRow
	- BatchWriteRow
	- GetRange
	- XGetRange


## Install

	$ go get code.google.com/p/goprotobuf/{proto,protoc-gen-go}
	$ go get github.com/GiterLab/goots
> **NOTE**: If you can't get `goprotobuf` package (you known why)，Please refer to [gopm.io](http://gopm.io/download) to download manually.

## Usage
	// delete a table

	// create a table
	table_meta := &OTSTableMeta{
		TableName: "myTable",
		SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
			"gid": "INTEGER",
			"uid": "INTEGER",
		},
	}

	reserved_throughput := &OTSReservedThroughput{
		OTSCapacityUnit{100, 100},
	}

	ots_err := ots_client.CreateTable(table_meta, reserved_throughput)
	if ots_err != nil {
		fmt.Println(ots_err)
	}
	fmt.Println("表已创建")

	// list tables
	list_tables, ots_err := ots_client.ListTable()
	if ots_err != nil {
		fmt.Println(ots_err)
	}
	fmt.Println("表的列表如下：")
	fmt.Println("list_tables:", list_tables.TableNames)

	// insert a row

	// get a row

More examples, please see [example/interfaces.go](https://github.com/GiterLab/goots/blob/master/example/interfaces.go).

## Links
- [Open Table Service，OTS](http://www.aliyun.com/product/ots)
- [OTS介绍](http://help.aliyun.com/list/11115779.html?spm=5176.383723.9.2.RYJAsQ)
- [OTS产品文档](http://oss.aliyuncs.com/aliyun_portal_storage/help/ots/OTS%20User%20Guide_Protobuf%20API%202%200%20Reference.pdf?spm=5176.383723.9.7.RYJAsQ&file=OTS%20User%20Guide_Protobuf%20API%202%200%20Reference.pdf)
- [使用API开发指南](http://help.aliyun.com/view/11108328_13761831.html?spm=5176.383723.9.6.RYJAsQ)
- [Python SDK开发包](http://oss.aliyuncs.com/aliyun_portal_storage/help/ots/ots_python_sdk_2.0.2.zip?spm=5176.383723.9.8.RYJAsQ&file=ots_python_sdk_2.0.2.zip)
- [Java SDK开发包](http://oss.aliyuncs.com/aliyun_portal_storage/help/ots/aliyun-openservices-OTS-2.0.4.zip?spm=5176.383723.9.9.RYJAsQ&file=aliyun-openservices-OTS-2.0.4.zip)
- [nodejs SDK](https://github.com/alibaba/ots)

## License

This project is under the MIT License. See the [LICENSE](https://github.com/GiterLab/goots/blob/master/LICENSE) file for the full license text.