#### 何时使用

弹出一个下拉菜单给用户选择操作

#### 基本使用

```jsx
initialState = {
    value:1, 
    data:[
        {value:1,label:'选项一',disabled:false},
        {value:2,label:'默认显示六项',disabled:false},
        {value:3,label:'不出现滚动条',disabled:false},
        {value:4,label:'显示七项时，出现纵向滚动条',disabled:false},
        {value:5,label:'宽度和高度可自定义',disabled:false},
        {value:6,label:'选项超过长度时候截断显示这是超过的部分超过的部分',disabled:false},
    ], 
};
<Select2
    value={state.value}
    onChange={(event) => setState({value: event.detail})}
>
    {
        state.data.map(item=>{
            return <Select2.Option key={item.value} value={item.value} disabled={item.disabled}>{item.label}</Select2.Option>
        })
    }
</Select2>
```

#### 下拉选项中带有图标

```jsx
initialState = {
    value:1, 
    data:[
        {value:1,label:'选项一',disabled:false},
        {value:2,label:'选项二',disabled:false},
        {value:3,label:'选项三',disabled:false},
    ] 
};
<Select2
    value={state.value}
    onChange={(event) => setState({value: event.detail})}
>
    {
        state.data.map(item=>{
            return <Select2.Option key={item.value} value={item.value} iconName={'filter'} disabled={item.disabled}>{item.label}</Select2.Option>
        })
    }
</Select2>
```

#### 默认超过6项出现滚动条（最大高度可控）

```jsx
initialState = {
    value:1, 
    data:[
        {value:1,label:'选项一',disabled:false},
        {value:2,label:'选项二',disabled:false},
        {value:3,label:'选项三',disabled:false},
        {value:4,label:'选项四',disabled:false},
        {value:5,label:'选项五',disabled:false},
        {value:6,label:'选项六',disabled:false},
        {value:7,label:'选项七',disabled:false},
    ], 
};
<Select2
    value={state.value}
    onChange={(event) => setState({value: event.detail})}
>
    {
        state.data.map(item=>{
            return <Select2.Option key={item.value} value={item.value} disabled={item.disabled}>{item.label}</Select2.Option>
        })
    }
</Select2>
```

#### 错误样式

```jsx
initialState = {
    value:1, 
    data:[
        {value:1,label:'选项一',disabled:false},
        {value:2,label:'选项二',disabled:false},
        {value:3,label:'选项三',disabled:false},
    ],
};

<Select2
    value={state.value}
    placeholder={'错误状态'}
    status={'error'}
    onChange={(event) => setState({value: event.detail})}
>
    {
        state.data.map(item=>{
            return <Select2.Option key={item.value} value={item.value} disabled={item.disabled}>{item.label}</Select2.Option>
        })
    }
</Select2>
```

#### 禁用-禁用时显示上一次选择值

```jsx
const Button = require('../Button').default;
initialState = {
    value:1, 
    disabled:false,
    data:[
        {value:1,label:'选项一',disabled:false},
        {value:2,label:'选项二',disabled:false},
        {value:3,label:'选项三',disabled:false},
    ] 
};

<div>
    <Select2
        value={state.value}
        disabled={state.disabled}
        onChange={(event) => setState({value: event.detail})}
    >
        {
            state.data.map(item=>{
                return <Select2.Option key={item.value} value={item.value} disabled={item.disabled}>{item.label}</Select2.Option>
            })
        }
    </Select2>
    <div style={{marginTop:'10px'}}>
        <Button onClick={()=> setState({disabled:!state.disabled})}>{state.disabled?'启用':'禁用'}</Button>
    </div>
</div>
```

#### 禁用-禁用时显示空

```jsx
const Button = require('../Button').default;
initialState = {
    value:1, 
    lastSelected:1,
    disabled:false,
    data:[
        {value:1,label:'选项一',disabled:false},
        {value:2,label:'选项二',disabled:false},
        {value:3,label:'选项三',disabled:false}
    ], 
};

<div>
    <Select2
        value={state.value}
        disabled={state.disabled}
        onChange={(event) => {
            setState({value: event.detail,lastSelected:event.detail})
        }}
    >
        {
            state.data.map(item=>{
                return <Select2.Option key={item.value} value={item.value} disabled={item.disabled}>{item.label}</Select2.Option>
            })
        }
    </Select2>
    <div style={{marginTop:'10px'}}>
        <Button onClick={()=> {
            setState({
                disabled:!state.disabled,
                value:state.disabled ?state.lastSelected:null
            })
        }}>{state.disabled?'启用':'禁用'}</Button>
    </div>
</div>
```

#### 禁用-禁用时显示placeholder

```jsx
const Button = require('../Button').default;

class Select2Demo extends React.Component {
     constructor() {
        super()
        this.state = {
            value:1,
            disabled:false,
            data:[{value:1,label:'选项一',disabled:false},{value:2,label:'选项二',disabled:false},{value:3,label:'选项三',disabled:false}],
            lastSelected:1
        }
    }

    render(){
        return (
            <div>
                <Select2
                    value={this.state.value}
                    disabled={this.state.disabled}
                    placeholder={'请选择一项'}
                    onChange={(event) => {
                        this.setState({value: event.detail,lastSelected:event.detail})
                    }}
                >
                    {
                        this.state.data.map(item=>{
                            return <Select2.Option key={item.value} value={item.value} disabled={item.disabled}>{item.label}</Select2.Option>
                        })
                    }
                </Select2>
                <div style={{marginTop:'10px'}}>
                    <Button onClick={()=> {
                        this.setState({
                            disabled:!this.state.disabled,
                            value:this.state.disabled ? this.state.lastSelected:null
                    })}}>
                        {this.state.disabled?'启用':'禁用'}
                    </Button>
                </div>
            </div>
        )}
};<Select2Demo/>

```