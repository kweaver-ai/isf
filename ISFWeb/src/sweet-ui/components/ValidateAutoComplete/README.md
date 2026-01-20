### 这个控件叫什么

自动完成验证框

### 何时使用

* 需要对自动完成下拉输入框输入的内容进行校验时

### 使用示例

#### 1.输入发生变化时验证

* 输入过程中触发验证。
* 当输入的内容不符合验证规则时，输入框显示为错误样式。当鼠标聚焦或者悬浮到输入框的时候，气泡提示错误详情。当再次输入的内容符合验证规则时候，输入框错误样式消失。

```jsx
const React = require('react');
const {range} = require('lodash');

const data = range(0,20).map((i) => ({id: i, text: 'xxxx'+i}))

class ValidateAutoCompleteLazyLoader extends React.Component {
    constructor() {
        super()
        this.state = {
            value: '',
            data: data.slice(0, 5),
            selection: null,
            validateState: 'normal',
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
            selection: (this.state.selection && this.state.selection.text === detail)? 
                this.state.selection : null,
            validateState: 'normal',
        })
    }

    select({detail}) {
        if(detail) {
            this.setState({
                selection: detail, 
                value: detail.text,
                validateState: 'invalid',
            })
        }
    }

    render() {
        const {value, data, selection, validateState} = this.state

        return (
            <div>
                <ValidateAutoComplete
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
                    validateState={validateState}
                    validateMessages={{
                            ['empty']: '输入不允许为空。',
                            ['invalid']: '该选项不存在。',
                    }}
                    loader={this.loader}
                    onLoad={this.loadData}
                    ListEmptyComponent={() => this.state.value ? <div>暂无匹配结果</div> : <div>无可选项，请添加</div>}
                    onValueChange={this.changeValue}
                    dropdownRender={() => (
                        <Button
                            key={'btn'}
                            width={'100%'}
                            style={{borderRadius: 0, border: 0, borderTop: '1px solid #cfcfcf'}}
                            onClick={() => alert('添加')}
                        >
                            {'添加'}
                        </Button>
                    )}
                    onSelect={this.select}
                    onPressEnter={this.select}
                />
                <Button
                    key={'btn'}
                    width={100}
                    style={{marginLeft: 20}}
                    onClick={() => this.setState({validateState: selection ? 'invalid' : value ? 'normal' : 'empty'})}
                >
                    {'确定'}
                </Button>
            </div>
        )
    }
};<ValidateAutoCompleteLazyLoader />
```