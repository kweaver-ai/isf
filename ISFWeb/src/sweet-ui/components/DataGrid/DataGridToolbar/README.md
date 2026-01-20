## DataGridToolbar

## Usage

```react
const ToolbarComponent = ({ data, selection, }) => {
    return (
        <div>
            <Button
                style={{
                    marginRight: 10
                }}
            >
                {"复制"}
            </Button>

            <Button
                style={{
                    marginRight: 10
                }}
            >
                {"新建文件夹"}
            </Button>
        </div>
    )
}
const data = [
    {
        name: 'Foo',
        age: 1
    },
    {
        name: 'Bar',
        age: 20
    },
]
<DataGridToolbar
    data={data}
    selection={[]}
    selectAll={(checked)=>{}}
    ToolbarComponent={ToolbarComponent}
/>
```

## Props

| 属性                                  | 说明                                     | 必须 | 类型            | 默认值 |
| ------------------------------------- | ---------------------------------------- | ---- | --------------- | ------ |
| data                                  | 表格数据源                               | No   | array           | []     |
| selection                             | 选中的数据项                             | No   | any             | /      |
| [selectAll](#selectAll)               | 全选按钮勾选状态发生改变时触发的处理函数 | No   | function        | /      |
| [ToolbarComponent](#ToolbarComponent) | 额外的工具栏组件                         | No   | React.ReactNode | /      |

## Reference

### Props

#### <span id="selectAll">selectAll</span>

全选按钮勾选状态发生改变时触发的处理函数。

| 参数    | 说明                     | 类型    |
| ------- | ------------------------ | ------- |
| checked | 全选按钮改变后的勾选状态 | boolean |

#### <span id="ToolbarComponent">ToolbarComponent</span>

额外的工具栏组件。

| 属性      | 说明         | 类型  |
| --------- | ------------ | ----- |
| data      | 表格数据源   | array |
| selection | 选中的数据项 | any   |
