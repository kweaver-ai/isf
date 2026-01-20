# DataTable

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
            },
            {
                name: 'Baz',
                age: 20,
                address: {
                    country: 'JP'
                }
            },
            {
                name: 'Faz',
                age: 21,
                address: {
                    country: 'KR'
                }
            }
        ]
const ToolbarComponent = ({ data, selection, selectAll }) => <div>render ToolbarComponent in here</div>
const HeaderComponent = ({ columns }) => <div>render HeaderComponent in here</div>
const FooterComponent = ({ data }) => <div>render FooterComponent in here</div>
const EmptyComponent = () => <div>render EmptyComponent in here</div>
<DataTable
    columns={columns}
    data={data}
    rowKeyExtractor={(record) => ''}
    cellKeyExtractor={(record, key) => ''}
    enableMultiSelect={false}
    selection={[data[0], data[2]]}
    ToolbarComponent={ToolbarComponent}
    HeaderComponent={HeaderComponent}
    FooterComponent={FooterComponent}
    EmptyComponent={EmptyComponent}
    onSelectionChange={selection=>{}}
/>
```

## Props

| 属性                                    | 说明                           | 必须 | 类型            | 默认值 |
| --------------------------------------- | ------------------------------ | ---- | --------------- | ------ |
| [columns](#columns)                     | 数据列配置                     | Yes  | array           | /      |
| data                                    | 表格数据源                     | No   | array           | /      |
| [rowKeyExtractor](#rowKeyExtractor)     | 生成每一行对应的 key           | No   | function        | /      |
| [cellKeyExtractor](#cellKeyExtractor)   | 生成每一个单元格对应的 key     | No   | function        | /      |
| enableMultiSelect                       | 是否允许多选                   | No   | boolean         | /      |
| [selection](#selection)                 | 默认选中数据                   | No   | any             | /      |
| [ToolbarComponent](#ToolbarComponent)   | 工具栏组件                     | No   | React.ReactNode | /      |
| [HeaderComponent](#HeaderComponent)     | 列表头组件                     | No   | React.ReactNode | /      |
| [FooterComponent](#FooterComponent)     | 列表页脚组件                   | No   | React.ReactNode | /      |
| [EmptyComponent](#EmptyComponent)       | 数据为空时显示的组件           | No   | React.ReactNode | /      |
| [onSelectionChange](#onSelectionChange) | 选中项发生变化时触发的处理函数 | No   | function        | /      |

## Reference

### Props

#### <span id="columns">columns</span>

列配置对象数组，每一项列配置对象包含以下属性：

| 属性       | 说明                                                                                             | 必须 | 类型               |
| ---------- | ------------------------------------------------------------------------------------------------ | ---- | ------------------ |
| width      | 列宽度，多个列会按比例分配，建议使用百分比                                                       | No   | string&#124;number |
| key        | 用于从多层级的对象中查找特定属性，会将结果传递给 renderCell(value, record)                       | No   | string             |
| renderCell | 渲染单元格。`value` 会通过 `key` 属性在 `record` 对象中进行查找。`record` 是该行完整的数据对象。 | No   | function           |
| title      | 列标题                                                                                           | No   | string             |

#### <span id="rowKeyExtractor">rowKeyExtractor</span>

生成该行数据对应的 `key` 属性。

| 参数   | 说明               | 类型 |
| ------ | ------------------ | ---- |
| record | 该行完整的数据对象 | any  |

#### <span id="cellKeyExtractor">cellKeyExtractor</span>

生成单元格对应的 `key` 属性。

| 参数   | 说明                             | 类型   |
| ------ | -------------------------------- | ------ |
| record | 该行完整的数据对象               | any    |
| key    | 该列 column 配置对象中对应的 key | string |

#### <span id="selection">selection</span>

默认选中数据项，selection 属性必须是 data 的子集或者其中一项。

> slection 必须保证和 data 引用一致。

```jsx
// 不会选中，因为selection和data中的第一项不指向同一对象
const data = [{name: 1}]
<DataTable data={data} selection={{name: 1}} />

// 会选中，因为selection === data
const data = [{name: 1}]
<DataTable data={data} selection={data[0]} />
```

#### <span id="ToolbarComponent">ToolbarComponent</span>

工具栏组件

| 参数      | 说明                                           | 类型     |
| --------- | ---------------------------------------------- | -------- |
| data      | 列表中的所有数据                               | array    |
| selection | 当前选中项                                     | array    |
| selectAll | (selected: boolean) => void;选中或去选中所有项 | function |

#### <span id="HeaderComponent">HeaderComponent</span>

表头组件

| 参数    | 说明       | 类型  |
| ------- | ---------- | ----- |
| columns | 数据列配置 | array |

#### <span id="FooterComponent">FooterComponent</span>

页脚组件

| 参数 | 说明             | 类型  |
| ---- | ---------------- | ----- |
| data | 列表中的所有数据 | array |

#### <span id="EmptyComponent">EmptyComponent</span>

内容为空提示组件

#### <span id="onSelectionChange">onSelectionChange</span>

选中项发生变化时触发的处理函数。

| 参数      | 说明                                                   | 类型 |
| --------- | ------------------------------------------------------ | ---- |
| selection | 单选时为选中的数据对象，多选时为选中数据对象构成的数组 | any  |
