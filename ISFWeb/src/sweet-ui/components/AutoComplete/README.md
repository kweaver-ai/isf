# AutoComplete

自动完成控件

## 基本使用

1. 搜索下拉
```jsx
const React = require('react');
const {range} = require('lodash');

const data = range(0,20).map((i) => ({id: i, text: 'xxxx'+i}))

class AutoCompleteLazyLoader extends React.Component {
    constructor() {
        super()
        this.state = {
            value: '',
            data: [],
            selection: null,
        }

        this.isLazyLoad = false

        this.loader = this.loader.bind(this)
        this.loadData = this.loadData.bind(this)
        this.changeValue = this.changeValue.bind(this)
        this.select = this.select.bind(this)
    }

    loader({key, start, limit}) {
        this.isLazyLoad = !!start

        if(key) {
            const all = data.filter((item) => item.text.includes(key))
            return all.slice(start, start + limit)
        }

        return []
    }

    loadData({detail}) {
        this.setState({data: this.isLazyLoad ? [...this.state.data, ...detail] : detail})
    }

    changeValue({detail}) {
        this.setState({
            value: detail,
            selection: (this.state.selection && this.state.selection.text === detail)? this.state.selection : null
        })
    }

    select({detail}) {
        if(detail) {
            this.setState({selection: detail, value: detail.text})
        }
    }

    render() {
        const {value, data, selection} = this.state
        console.log(value)

        return (
            <div style={{ display: 'inline-block', width: '200px'}}>
                <AutoComplete
                    maxHeight={100}
                    width={280}
                    selection={selection}
                    value={value}
                    placeholder={'搜索'}
                    data={data}
                    limit={5}
                    allowClear={true}
                    renderItem={(item, index) => (`row${index} ${item.text}`)}
                    loader={this.loader}
                    onLoad={this.loadData}
                    onValueChange={this.changeValue}
                    onSelect={this.select}
                    onPressEnter={this.select}
                    {...!this.state.data.length && this.state.value ? {ListEmptyComponent: <div>暂无匹配结果</div>}: {}}
                />
            </div>
        )
    }
};<AutoCompleteLazyLoader />

```


2. 下拉搜索选择，可新建
```jsx
const React = require('react');
const {range} = require('lodash');

const data = range(0,20).map((i) => ({id: i, text: 'xxxx'+i}))

class AutoCompleteLazyLoader extends React.Component {
    constructor() {
        super()
        this.state = {
            value: '',
            data: data.slice(0, 5),
            selection: null,
        }

        this.isLazyLoad = false

        this.loader = this.loader.bind(this)
        this.loadData = this.loadData.bind(this)
        this.changeValue = this.changeValue.bind(this)
        this.select = this.select.bind(this)
    }

    loader({key, start, limit}) {
        this.isLazyLoad = !!start

        if(key) {
            const all = data.filter((item) => item.text.includes(key))
            return all.slice(start, start + limit)
        }

        return data.slice(start, start + limit)
    }

    loadData({detail}) {
        this.setState({data: this.isLazyLoad ? [...this.state.data, ...detail] : detail})
    }

    changeValue({detail}) {
        this.setState({
            value: detail,
            selection: (this.state.selection && this.state.selection.text === detail)? this.state.selection : null
        })
    }

    select({detail}) {
        if(detail) {
            this.setState({selection: detail.selection, value: detail.text})
        }
    }

    render() {
        const {value, data, selection} = this.state

        return (
            <div style={{ display: 'inline-block', width: '200px'}}>
                <AutoComplete
                    maxHeight={100}
                    width={280}
                    selection={selection}
                    value={value}
                    placeholder={'请选择'}
                    data={data}
                    limit={5}
                    allowClear={true}
                    iconOnBefore={''}
                    iconOnAfter={'arrowDown'}
                    renderItem={(item, index) => (`row${index} ${item.text}（过过过过过过过过过过过过过过过过过）`)}
                    loader={this.loader}
                    onLoad={this.loadData}
                    ListEmptyComponent={() => this.state.value ? <div>暂无匹配结果</div> : <div>无可选项，请添加</div>}
                    onValueChange={this.changeValue}
                    dropdownRender={() => (
                        <Button
                            key={'btn'}
                            width={'100%'}
                            style={{borderRadius: 0, border: 0, borderTop: '1px solid #0000000d'}}
                            onClick={() => alert('添加')}
                        >
                            {'添加'}
                        </Button>
                    )}
                    onSelect={this.select}
                    onPressEnter={this.select}
                />
            </div>
        )
    }
};<AutoCompleteLazyLoader />

```

3. 多选
```jsx
const React = require('react');
const {range} = require('lodash');

const data = range(0,20).map((i) => ({id: i, text: 'xxxx'+i}))

class AutoCompleteLazyLoader extends React.Component {
    constructor() {
        super()
        this.state = {
            value: '',
            data: data.slice(0, 5),
            selection: [],
        }

        this.isLazyLoad = false

        this.loader = this.loader.bind(this)
        this.loadData = this.loadData.bind(this)
        this.changeValue = this.changeValue.bind(this)
        this.select = this.select.bind(this)
    }

    loader({key, start, limit}) {
        this.isLazyLoad = !!start

        if(key) {
            const all = data.filter((item) => item.text.includes(key))
            return all.slice(start, start + limit)
        }

        return data.slice(start, start + limit)
    }

    loadData({detail}) {
        this.setState({data: this.isLazyLoad ? [...this.state.data, ...detail] : detail})
    }

    changeValue({detail}) {
        this.setState({
            value: detail,
        })
    }

    select({detail}) {
        console.log(detail)
        this.setState({
            selection: detail && detail.selection,
        })
    }

    render() {
        const {value, data, selection} = this.state

        return (
            <div>
                <AutoComplete
                    maxHeight={120}
                    width={280}
                    selection={selection}
                    value={value}
                    placeholder={'请选择'}
                    data={data}
                    limit={5}
                    allowClear={true}
                    enableMultiSelect={true}
                    iconOnBefore={''}
                    iconOnAfter={'arrowDown'}
                    renderItem={(item, index) => (`row${index} ${item.text}（过过过过过过过过过过过过过过过过过）`)}
                    loader={this.loader}
                    onLoad={this.loadData}
                    ListEmptyComponent={ <div>暂无匹配结果</div>}
                    onValueChange={this.changeValue}
                    onSelect={this.select}
                    onPressEnter={this.select}
                />
                <div style={{display: 'inline-block', marginLeft: 20, width: 300}}>
                    已选：{selection.map(({text}) => text).join(',')}
                </div>
            </div>
        )
    }
};<AutoCompleteLazyLoader />

```

4. 多选非懒加载 —— 可全选
```jsx
const React = require('react');
const {range} = require('lodash');

const data = range(0,20).map((i) => ({id: i, text: 'xxxx'+i}))

class AutoCompleteLazyLoader extends React.Component {
    constructor() {
        super()
        this.state = {
            value: '',
            data: [],
            selection: [],
        }

        this.loader = this.loader.bind(this)
        this.loadData = this.loadData.bind(this)
        this.changeValue = this.changeValue.bind(this)
        this.select = this.select.bind(this)
    }

    loader({key}) {

        if(key) {
            return data.filter((item) => item.text.includes(key))
        }

        return data
    }

    loadData({detail}) {
        this.setState({data: detail})
    }

    changeValue({detail}) {
        this.setState({
            value: detail,
        })
    }

    select({detail}) {
        console.log(detail)
        if(detail) {
            const {selection, isSelectedAll} = detail

            this.setState({
                selection,
            })
        }
    }

    render() {
        const {value, data, selection} = this.state

        return (
            <div>
                <AutoComplete
                    maxHeight={120}
                    width={280}
                    selection={selection}
                    value={value}
                    placeholder={'搜索'}
                    data={data}
                    enableMultiSelect={true}
                    enableSelectAll={true}
                    renderItem={(item, index) => (`row${index} ${item.text}`)}
                    loader={this.loader}
                    onLoad={this.loadData}
                    ListEmptyComponent={ <div>暂无匹配结果</div>}
                    onValueChange={this.changeValue}
                    onSelect={this.select}
                    onPressEnter={this.select}
                />
                <div style={{display: 'inline-block', marginLeft: 20, width: 300}}>
                    已选：{selection.map(({text}) => text).join(',')}
                </div>
            </div>
        )
    }
};<AutoCompleteLazyLoader />

```

5. 多选懒加载 —— 可全选
```jsx
const React = require('react');
const {range} = require('lodash');

const data = range(0,20).map((i) => ({id: i, text: 'xxxx'+i}))

class AutoCompleteLazyLoader extends React.Component {
    constructor() {
        super()
        this.state = {
            value: '',
            data: data.slice(0, 5),
            selection: [],
        }

        this.isLazyLoad = false

        this.loader = this.loader.bind(this)
        this.loadData = this.loadData.bind(this)
        this.changeValue = this.changeValue.bind(this)
        this.select = this.select.bind(this)
    }

    loader({key, start, limit}) {
        this.isLazyLoad = !!start

        if(key) {
            const all = data.filter((item) => item.text.includes(key))
            return all.slice(start, start + limit)
        }

        return data.slice(start, start + limit)
    }

    loadData({detail}) {
        this.setState({data: this.isLazyLoad ? [...this.state.data, ...detail] : detail})
    }

    changeValue({detail}) {
        this.setState({
            value: detail,
        })
    }

    select({detail}) {
        console.log(detail)
        if(detail) {
            const {selection, isSelectedAll} = detail

            // 搜索的结果
            const res = data.filter((item) => item.text.includes(this.state.value))

            this.setState({
                selection: isSelectedAll ?
                    [
                        ...this.state.selection, 
                        ...res.filter((item) => !this.state.selection.includes(item))
                    ] 
                    : selection,
            })
        }
    }

    render() {
        const {value, data, selection} = this.state

        return (
            <div>
                <AutoComplete
                    maxHeight={120}
                    width={280}
                    selection={selection}
                    value={value}
                    placeholder={'搜索'}
                    data={data}
                    limit={5}
                    enableMultiSelect={true}
                    enableSelectAll={true}
                    renderItem={(item, index) => (`row${index} ${item.text}`)}
                    loader={this.loader}
                    onLoad={this.loadData}
                    ListEmptyComponent={ <div>暂无匹配结果</div>}
                    onValueChange={this.changeValue}
                    onSelect={this.select}
                    onPressEnter={this.select}
                />
                <div style={{display: 'inline-block', marginLeft: 20, width: 300}}>
                    已选：{selection.map(({text}) => text).join(',')}
                </div>
            </div>
        )
    }
};<AutoCompleteLazyLoader />

```