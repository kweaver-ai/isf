# BaseTable

基础的表格组件。

## Usage

```react
const columns = [
            {
                width: '30%',
                key: 'name',
                renderCell: (name, o) => name,
            },
            {
                width: '20%',
                key: 'age',
                renderCell: (age, o) => age,
            },
            {
                width: '50%',
                key: 'address.country',
                renderCell: (country, o) => country,
            }
        ]
const data = [
    {
        name: 'Foo',
        age: 30,
        address: {
            country: 'China'
        }
    },
    {
        name: 'Bar',
        age: 29,
        address: {
            country: 'US'
        }
    }
]
<BaseTable
    columns={columns}
    data={data}
    rowKeyExtractor={(record) => ''}
    cellKeyExtractor={(record, key) => ''}
    rowClassName={(record, index) => ''}
    cellClassName={(record, key, index) => ''}
    rowHoverClassName="customizeRowHoverClassName"
    onRowClicked={(record,index,event)=>{}}
    onRowEnter={(record,index)=>{}}
    onRowLeave={(record,index)=>{}}
/>
```

## Props

| 参数                                    | 说明                         | 必须 | 类型                 | 默认值 |
| :-------------------------------------- | :--------------------------- | :--- | :------------------- | :----- |
| [columns](#columns)                     | 数据列配置                   | Yes  | array                | /      |
| data                                    | 表格数据源                   | No   | array                | /      |
| [rowClassName](#rowClassName)           | 每一行 className             | No   | string&#124;function | /      |
| [cellClassName](#cellClassName)         | 每一个单元格 className       | No   | string&#124;function | /      |
| [rowHoverClassName](#rowHoverClassName) | 行在鼠标悬浮时的 className   | No   | string               | /      |
| [rowKeyExtractor](#rowKeyExtractor)     | 生成每一行对应的 key         | No   | function             | /      |
| [cellKeyExtractor](#cellKeyExtractor)   | 生成每一个单元格对应的 key   | No   | function             | /      |
| [onRowEnter](#onRowEnter)               | 鼠标移入一行时触发的处理函数 | No   | function             | /      |
| [onRowLeave](#onRowLeave)               | 鼠标移出一行时触发的处理函数 | No   | function             | /      |
| [onRowClicked](#onRowClicked)           | 鼠标单击一行时触发的处理函数 | No   | function             | /      |

## Reference

### Props

#### <span id="columns">columns</span>

列配置对象数组，每一项列配置对象包含以下属性：

| 属性       | 说明                                                                                             | 必须 | 类型               |
| :--------- | :----------------------------------------------------------------------------------------------- | ---- | :----------------- |
| width      | 列宽度，多个列会按比例分配，建议使用百分比                                                       | No   | string&#124;number |
| key        | 用于从多层级的对象中查找特定属性，会将结果传递给 renderCell(value, record)                       | No   | string             |
| renderCell | 渲染单元格。`value` 会通过 `key` 属性在 `record` 对象中进行查找。`record` 是该行完整的数据对象。 | No   | function           |

Example:

```jsx
columns: [
    {
        key: 'name',
        width: '30%',
        renderCell: (name, record) => name,
    },
    {
        key: 'detail.description',
        width: '70%',
        renderCell: (description, record) => description,
    },
];
```

#### <span id="rowClassName">rowClassName</span>

每一行的 className。

可选传入类型：

*   sring：直接应用到所有列表行
*   function：根据参数生成返回 className,具体参数见下表。

|  参数  |          说明          | 类型   |
| :----: | :--------------------: | ------ |
| record |  当前行完整的数据对象  | any    |
| index  | 当前行在数据表中的索引 | number |

#### <span id="cellClassName">cellClassName</span>

每个单元格的 className。

可选传入类型：

*   sring：直接应用到所有列表单元格。
*   function：根据参数生成返回 className,具体参数见下表。

|  参数  |               说明               | 类型   |
| :----: | :------------------------------: | ------ |
| record |        该行完整的数据对象        | any    |
|  key   | 该列 column 配置对象中对应的 key | string |
| index  |   当前单元格在当前行中的索引值   | number |

#### <span id="rowHoverClassName">rowHoverClassName</span>

悬浮时每一行的 className。

可选传入类型：

*   sring：直接应用到每个列表行。
*   function：根据参数生成返回 className,具体参数见下表。

| 参数   | 说明                     | 类型   |
| ------ | ------------------------ | ------ |
| record | 当前行完整的数据对象     | any    |
| index  | 当前行在数据表中的索引值 | number |

#### <span id="rowKeyExtractor">rowKeyExtractor</span>

生成该行数据对应的 `key` 属性。

|  参数  |        说明        | 类型 |
| :----: | :----------------: | ---- |
| record | 该行完整的数据对象 | any  |

#### <span id="cellKeyExtractor">cellKeyExtractor</span>

生成单元格对应的 `key` 属性。

|  参数  |               说明               |  类型  |
| :----: | :------------------------------: | :----: |
| record |        该行完整的数据对象        |  any   |
|  key   | 该列 column 配置对象中对应的 key | string |

#### <span id="onRowEnter">onRowEnter</span>

鼠标进入某一行时触发。

|  参数  |        说明        |  类型  |
| :----: | :----------------: | :----: |
| record | 该行完整的数据对象 |  any   |
| index  |  该行所在位置下标  | number |

#### <span id="onRowLeave">onRowLeave</span>

鼠标离开某一行时触发。

|  参数  |        说明        |  类型  |
| :----: | :----------------: | :----: |
| record | 该行完整的数据对象 |  any   |
| index  |  该行所在位置下标  | number |

#### <span id="onRowClicked">onRowClicked</span>

鼠标单击某一行时触发。

|  参数  |        说明        |                    类型                     |
| :----: | :----------------: | :-----------------------------------------: |
| record | 该行完整的数据对象 |                     any                     |
| index  |  该行所在位置下标  |                   number                    |
| event  |    鼠标事件对象    | React.MouseEvent&lt;HTMLTableRowElement&gt; |
