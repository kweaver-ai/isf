### 何时使用

* 当有大量结构化的数据需要展现时；

### 基本使用

> 指定`DataGrid`数据源`data`和列配置对象`columns`为数组。rowKeyExtractor传入用来计算每一行的key值的方法，cellKeyExtractor生成每行对应的唯一key。

```jsx
const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    }
];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('removed')}}>remove</button>,
    },
];

<DataGrid
    columns={columns}
    height={250}
    data={data}
    rowKeyExtractor={({ id }) => id}
    cellKeyExtractor={(record, key) => `${record.id}-${key}`}
/>

```

> 设置`headless`属性为`true`实现取消表头的效果。

```jsx
const React = require('react');

const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    }
];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('removed')}}>remove</button>,
    },
];

<DataGrid
    headless={true}
    columns={columns}
    height={250}
    data={data}
    rowKeyExtractor={({ id }) => id}
    cellKeyExtractor={(record, key) => `${record.id}-${key}`}
/>

```
#### 1.分页模式

* 指定分页配置参数`DataGridPager`，就变成了支持分页显示的表格控件。当前页码发生变化时触发`onPageChange`事件，传递参数为下一次的页码`page`和每页大小`size`。

```jsx
const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    },
    {
        name: 'Jack',
        age: 29,
        id: 5,
    },
    {
        name: 'Joy',
        age: 23,
        id: 6,
    },
    {
        name: 'Seulgi',
        age: 25,
        id: 7,
    },
    {
        name: 'Wendy',
        age: 25,
        id: 8,
    },

];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('handleClick')}}>remove</button>,
    },
];

const EmptyComponent = () => <div>EmptyComponent</div>;

class DataGridWithPager extends React.Component {
    constructor() {
        super()
        this.state = {
            start: 0,
            end: 4
        }

        this.handlePageChange = this.handlePageChange.bind(this)
    }

    handlePageChange(event) {
        const {detail} = event;

        this.setState({
            start: (detail.page-1) * 4,
            end: detail.page * 4
        })
    }

    render() {
        return (
            <DataGrid
                rowHoverClassName={'row-hovering'}
                columns={columns}
                height={280}
                data={data.slice(this.state.start, this.state.end)}
                rowKeyExtractor={({ id }) => id}
                cellKeyExtractor={(record, key) => `${record.id}-${key}`}
                DataGridPager={{
                    size: 4,
                    onPageChange: this.handlePageChange,
                    total: 8
                }}
                onRowDoubleClicked={() => alert('onRowDoubleClicked')}
            />
        )
    }

};<DataGridWithPager />

```


#### 2.排序

> 指定sort排序参数执行初始化排序，包括排序字段`key`和排序方式`type`

> 排序字段所在列的表头会有一个箭头图标指示，箭头向上指示升序，向下指示降序；对应`columns`设置`sortable: true`表示该列允许排序，点击图标触发`onRequestSort`排序处理事件，传递参数包括排序字段`key`和排序方式`type`

```jsx
const SortType = require('./DataGridHeader').SortType;

const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    },
    {
        name: 'Jack',
        age: 29,
        id: 5,
    },
    {
        name: 'Joy',
        age: 23,
        id: 6,
    },
    {
        name: 'Seulgi',
        age: 25,
        id: 7,
    },
    {
        name: 'Wendy',
        age: 25,
        id: 8,
    },

];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        sortable: true,
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('handleClick')}}>remove</button>,
    },
];

const EmptyComponent = () => <div>EmptyComponent</div>;

class DataGridWithSort extends React.Component {
    constructor() {
        super()
        this.state = {
            start: 0,
            end: 4,
            loading: false,
            sort: {key: 'age', type: SortType.ASC},
            dataSource: data
        }

        this.requestSort = this.requestSort.bind(this)
        this.loadData = this.loadData.bind(this)
    }

    componentDidMount() {
        this.loadData(this.state.sort.type)
    }

    loadData(type) {
        const sortData = data.sort((a, b) => {
            return type === SortType.ASC ? (a.age - b.age) : (b.age - a.age)
        })
        this.setState({
            dataSource: sortData
        })
    }

    requestSort({key, type}) {
        this.setState({
            sort: {key, type}           
        })
        this.loadData(type)
    }

    render() {
        return (
            <DataGrid
                columns={columns}
                height={200}
                data={this.state.dataSource}
                rowKeyExtractor={({ id }) => id}
                cellKeyExtractor={(record, key) => `${record.id}-${key}`}
                sort={this.state.sort}
                onRequestSort={this.requestSort.bind(this)}
            />
        )
    }
};<DataGridWithSort />

```

* 分页模式下的排序

> 点击排序图标会同时触发`onParamsChange`和 `onRequestSort`事件，推荐使用`onParamsChange`， 该事件同时传递分页和排序参数；
> 如果当前不在首页，点击排序后会自动回到首页。

```jsx
const SortType = require('./DataGridHeader').SortType;

const PAGESIZE = 4;

const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    },
    {
        name: 'Jack',
        age: 29,
        id: 5,
    },
    {
        name: 'Joy',
        age: 23,
        id: 6,
    },
    {
        name: 'Seulgi',
        age: 25,
        id: 7,
    },
    {
        name: 'Wendy',
        age: 25,
        id: 8,
    },

];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        sortable: true,
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('handleClick')}}>remove</button>,
    },
];

class DataGridWithPageSort extends React.Component {
    constructor() {
        super()
        this.state = {
            loading: false,
            sort: {key: 'age', type: SortType.ASC},
            dataSource: data
        }

        this.handleParamsChange = this.handleParamsChange.bind(this)
        this.handlePageChange = this.handlePageChange.bind(this)
        this.loadData = this.loadData.bind(this)
    };

    

    componentDidMount() {
        this.loadData(0, this.state.sort.type)
    }

    loadData(start, type) {
        const sortData = data.slice(start, start + PAGESIZE).sort((a, b) => {
            return type === SortType.ASC ? 
            (a.age - b.age)
            : (b.age - a.age)
        })
        this.setState({
            dataSource: sortData
        })
    }

    handleParamsChange({page, sort}) {
        console.log({page,sort})
        if(sort) {
            this.setState({sort});

        }
        this.loadData(
            page ? (page-1) * PAGESIZE : 0, 
            sort ? sort.type : this.state.sort.type
        )
    }

    handlePageChange(event) {
        const {detail} = event;

        this.setState({
            dataSource: data.slice((detail.page-1) * 4, detail.page * 4)
        })
    }

    render() {
        return (
            <DataGrid
                columns={columns}
                height={200}
                data={this.state.dataSource}
                rowKeyExtractor={({ id }) => id}
                cellKeyExtractor={(record, key) => `${record.id}-${key}`}
                refreshing={this.state.loading}
                RefreshingComponent={() => <div>Refresing</div>}
                sort={this.state.sort}
                onParamsChange={this.handleParamsChange}
                DataGridPager={{
                    size: 4,
                    page: 1,
                    total: 8
                }}
            />
        )
    }
};<DataGridWithPageSort />

```

#### 3.懒加载

> 指定每次懒加载数据的起始索引`start`和每次加载的最大数据条数`limit`（不需要指定分页配置参数`DataGridPager`）。当滚动到接近列表底部的位置时会自动触发下一次懒加载，抛出`onRequestLazyLoad`事件，传递参数包括下一次加载数据的起始索引`start`和每次加载的最大数据条数`limit`。

```jsx
const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    },
    {
        name: 'Jack',
        age: 29,
        id: 5,
    },
    {
        name: 'Joy',
        age: 23,
        id: 6,
    },
    {
        name: 'Seulgi',
        age: 25,
        id: 7,
    },
    {
        name: 'Wendy',
        age: 25,
        id: 8,
    },

];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('handleClick')}}>remove</button>,
    },
];

const EmptyComponent = () => <div>EmptyComponent</div>;

class DataGridWithLazyLoader extends React.Component {
    constructor() {
        super()
        this.state = {
            start: 0,
            end: 4,
            loading: false
        }

        let loadComplete = false

        this.handleLazyLoad = this.handleLazyLoad.bind(this)
    }

    handleLazyLoad(event) {
        if(!this.loadComplete) {
            const {detail} = event;
            const {start, limit} = detail;

            this.setState({              
                loading: true
            })
            this.loadComplete = true

            setTimeout(() => {
                this.setState({
                    loading: false,
                    end: start + limit
                })
            }, 1000);
        }     
    }

    render() {
        return (
            <DataGrid
                EmptyComponent={EmptyComponent}
                rowHoverClassName={'row-hovering'}
                columns={columns}
                height={200}
                data={data.slice(0, this.state.end)}
                rowKeyExtractor={({ id }) => id}
                cellKeyExtractor={(record, key) => `${record.id}-${key}`}
                refreshing={this.state.loading}
                RefreshingComponent={() => <div>Loading More...</div>}
                start={this.state.start}
                limit={4}
                onRequestLazyLoad={this.handleLazyLoad}
            />
        )
    }
};<DataGridWithLazyLoader />

```
#### 4.选中

> 设置`enableSelect`为`true`使列表支持选中，单击某一行触发选中事件`onSelectionChange`，传递参数为选中的数据对象；

```jsx
const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    },
    {
        name: 'Jack',
        age: 29,
        id: 5,
    },
];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('removed')}}>remove</button>,
    },
];

    <DataGrid
        columns={columns}
        height={250}
        data={data}
        rowKeyExtractor={({ id }) => id}
        cellKeyExtractor={(record, key) => `${record.id}-${key}`}
        enableSelect={true}
        onSelectionChange={() => console.log('onSelectionChange')}   
    />
```

> 设置`enableSelect`和`enableMultiSelect`为`true`使列表支持多选；点击checkbox或单击某一行触发选中事件`onSelectionChange`，传递参数为选中数据对象构成的数组；

> DataGridHeader表头配置`enableSelectAll: true`表示表头显示【全选】复选框，默认值为`false`，勾选全选checkbox触发全选事件`onSelectAllChange`，传递参数类型为`boolean`，表示是否全选。

```jsx
/**
 * 内容为空提示组件
 */
const EmptyComponent = () => <div>EmptyComponent</div>;

const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    },
    {
        name: 'Jack',
        age: 29,
        id: 5,
    },
];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('removed')}}>remove</button>,
    },
];

    <DataGrid
        EmptyComponent={EmptyComponent}
        rowHoverClassName={'row-hovering'}
        columns={columns}
        height={250}
        data={data}
        rowKeyExtractor={({ id }) => id}
        cellKeyExtractor={(record, key) => `${record.id}-${key}`}
        enableSelect={true}
        enableMultiSelect={true}
        refreshing={false}
        RefreshingComponent={() => <div>Refresing</div>}
        DataGridHeader={{enableSelectAll: true}}
        onSelectionChange={() => console.log('onSelectionChange')}
        onChangeSelectAll={() => console.log('onChangeSelectAll')}
        onRowDoubleClicked={() => alert('onRowDoubleClicked')}
    />

```
#### 5.工具栏

> 通过添加参数`ToolbarComponent`指定工具栏内容，工具栏配置`DataGridToolbar`参数通过指定`enableSelectAll: true`在工具栏显示【全选】工具，默认值`false`。

```jsx
const Button = require('../Button').default;

/**
 * 工具栏组件
 */
const ToolbarComponent = ({ data, selection }) => {
    return (
        <div style={{ display: 'table', tableLayout: 'fixed', width: '100%' }}>
            <div style={{ display: 'table-cell', verticalAlign: 'middle', }}>
                <Button
                    style={{
                        marginRight: 10,
                        display: selection.length === data.length ? 'inline-block' : 'none',
                    }}
                >
                    {"复制"}
                </Button>
                <Button
                    style={{
                        marginRight: 10,
                        display: selection.length ? 'inline-block' : 'none',
                    }}
                    size={'auto'}
                >
                    {"新建文件夹"}
                </Button>
            </div>
            <div style={{ display: 'table-cell', verticalAlign: 'middle', textAlign: 'right' }}>
                <input style={{ height: 30, boxSizing: 'border-box' }} type="text" name="" id="" />
            </div>
        </div >
    );
};

const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    },
    {
        name: 'Jack',
        age: 29,
        id: 5,
    },
];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('removed')}}>remove</button>,
    },
];

        
    <DataGrid
        ToolbarComponent={ToolbarComponent}
        rowHoverClassName={'row-hovering'}
        columns={columns}
        height={250}
        data={data}
        rowKeyExtractor={({ id }) => id}
        cellKeyExtractor={(record, key) => `${record.id}-${key}`}
        enableSelect={true}
        enableMultiSelect={true}
        DataGridToolbar={{enableSelectAll: true}}
        onSelectionChange={() => console.log('onSelectionChange')}
        onChangeSelectAll={() => console.log('onChangeSelectAll')}
        onRowDoubleClicked={() => alert('onRowDoubleClicked')}
    />;

```


#### 6.嵌套子表格

> 展示每行数据更详细的信息。
> 添加`RowExtraComponent`渲染行额外信息，`showRowExtraOf`通过传入渲染行的下标或该行的数据用于判断在哪一行下渲染`RowExtraComponent`。
> 只支持同时显示一行的额外信息，请保证`showRowExtraOf`只对一行是生效的。

```jsx
const EmptyComponent = () => <div>EmptyComponent</div>;

const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    }
];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('removed')}}>remove</button>,
    },
];

const extradata = [
    {
        country: 'CHN',
        zone: 'ShangHai',
        id: 1,
    },
    {
        country: 'KR',
        zone: 'Seoul',
        id: 2,
    },
];

const extracolumns = [
    {
        key: 'country',
        width: '50%',
        title: 'Country',
        renderCell: (country, o) => country,
    },
    {
        key: 'zone',
        width: '50%',
        title: 'Zone',
        renderCell: (zone, o) => zone,
    },
];

const RowExtraComponent = () => {
    return (
        <DataGrid
            data={extradata}
            headless={true}
            refreshing={false}
            columns={extracolumns}
        />
    )
}

<DataGrid
    columns={columns}
    height={250}
    data={data}
    rowKeyExtractor={({ id }) => id}
    cellKeyExtractor={(record, key) => `${record.id}-${key}`}
    showRowExtraOf={0}
    RowExtraComponent={RowExtraComponent}
    onRowDoubleClicked={() => alert('onRowDoubleClicked')}
/>

```

#### 7.筛选过滤

> 在columns指定filters字段设置某一列的过滤项，勾选过滤项关闭弹窗后触发`onRequestFilter`事件

```jsx

const data = [
    {
        name: 'Foo',
        age: 30,
        id: 1,
    },
    {
        name: 'Bar',
        age: 29,
        id: 2,
    },
    {
        name: 'Baz',
        age: 29,
        id: 3,
    },
    {
        name: 'Tom',
        age: 29,
        id: 4,
    },
    {
        name: 'Jack',
        age: 29,
        id: 5,
    },
    {
        name: 'Joy',
        age: 23,
        id: 6,
    },
    {
        name: 'Seulgi',
        age: 25,
        id: 7,
    },
    {
        name: 'Wendy',
        age: 25,
        id: 8,
    },

];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        filters: [
            {
                text: 'Joy',
                value: 'Joy'
            },
            {
                text: 'Wendy',
                value: 'Wendy'
            }
        ],
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <button onClick={() => {alert('handleClick')}}>remove</button>,
    },
];

const EmptyComponent = () => <div>EmptyComponent</div>;

class DataGridWithFilter extends React.Component {
    constructor() {
        super()
        this.state = {
            start: 0,
            end: 4,
            loading: false,
            dataSource: data,
            selections: []
        }

        this.requestFilter = this.requestFilter.bind(this)
        this.loadData = this.loadData.bind(this)
        this.handleSelectionChange = this.handleSelectionChange.bind(this)
    }

    componentDidMount() {
        this.loadData(data)
    }
    handleSelectionChange({detail}) {
        this.setState({
            selections: detail
        })

    }

    loadData(data) {
        this.setState({
            dataSource: data
        })
    }

    requestFilter(event) {
        const {detail} = event

        alert('requestFilter')
    }

    render() {
        return (
            <DataGrid
                columns={columns}
                height={200}
                data={this.state.dataSource}
                selection={this.state.selections}
                rowKeyExtractor={({ id }) => id}
                cellKeyExtractor={(record, key) => `${record.id}-${key}`}
                onRequestFilter={this.requestFilter.bind(this)}
                onSelectionChange={this.handleSelectionChange.bind(this)}
                enableSelect={true}
                enableMultiSelect={true}
            />
        )
    }
};<DataGridWithFilter />

```

#### 8.拖拽排序(删除)

> 添加属性`draggable: true`使表格支持拖拽排序，拖拽排序结束后触发`onDragReorder`事件，传递参数为排序后的数据。

* 该特性暂不支持嵌套子表格

```jsx
const Button = require('../Button').default;

const data = [
    {
        name: 'Irene',
        age: 28,
        id: 1,
    },
    {
        name: 'Seulgi',
        age: 25,
        id: 2,
    },
    {
        name: 'Wendy',
        age: 25,
        id: 3,
    },
    {
        name: 'Joy',
        age: 23,
        id: 4,
    },
    {
        name: 'Yeri',
        age: 20,
        id: 5,
    }
];
const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <Button onClick={() => {alert('removed')}}>remove</Button>,
    },
];

initialState={data: data};

<DataGrid
    columns={columns}
    height={300}
    data={state.data}
    rowKeyExtractor={({ id }) => id}
    cellKeyExtractor={(record, key) => `${record.id}-${key}`}
    draggable={true}
    onDragReorder={({detail}) => {console.log(detail); setState({data: detail})} }
/>

```
#### 9.树形展开

> 添加属性`expandable: true`使表格支持树形展开，`rowKeyName`属性表示行唯一标识的名称，默认为`id`，`childrenColumnName`属性表示孩子纵队名称，默认为`children`，`expandedKeys: []`属性表示展开行的keys，点击展开/收起触发`onExpand`事件，传递参数为`expandedKeys`，`record`，`expanded`。

* 该特性暂不支持多选列表

```jsx
const data = [
    {
        name: 'Foo',
        age: 30,
        key: 1,
        children: [
            {
                name: 'Tom',
                age: 29,
                key: 11,
                children: [
                    {
                        name: 'Jon',
                        age: 2,
                        key: 111,
                    }
                ]
            }
        ]
    },
    {
        name: 'Bar',
        age: 22,
        key: 2,
    },
    {
        name: 'Lily',
        age: 12,
        key: 3,
        children: [
            {
                name: 'Kim',
                age: 29,
                key: 31,
                children: [
                    {
                        name: 'Jony',
                        age: 34,
                        key: 331,
                    }
                ]
            },
            {
                name: 'Hilen',
                age: 23,
                key: 32
            }
        ]
    }
];

const columns = [
    {
        key: 'name',
        width: '30%',
        title: 'Name',
        renderCell: (name, o) => name,
    },
    {
        key: 'age',
        width: '50%',
        title: 'Age',
        renderCell: (age, o) => age,
    },
    {
        width: '20%',
        title: 'Action',
        renderCell: () => <Button onClick={() => {alert('removed')}}>{'remove'}</Button>,
    },
];

class ExpandDataGrid extends React.Component {
    constructor() {
        super()
        this.state = {
            expandedKeys: [],
        }

        this.expand = this.expand.bind(this)
        this.reset = this.reset.bind(this)
    }

    expand({expandedKeys, record, expanded}) {
        this.setState({
            expandedKeys
        })
    }

    reset(){
        this.setState({
            expandedKeys: [],
        })
    }

    render() {
        return (
            <div>
                <Button
                    style={{
                        marginRight: 10,
                        display: 'inline-block',
                        marginBottom: 10
                    }}
                    width={120}
                    onClick={this.reset}
                >
                    {'收起所有展开'}
                </Button>
                <DataGrid
                    data={data}
                    columns={columns}
                    height={250}
                    enableSelect={true}
                    expandable={true}
                    rowKeyName={'key'}
                    childrenColumnName={'children'}
                    expandedKeys={this.state.expandedKeys}
                    onExpand={({detail}) => this.expand(detail)}
                />
            </div>
        )
    }
}

<ExpandDataGrid/>

```