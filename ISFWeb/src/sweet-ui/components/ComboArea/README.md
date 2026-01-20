#### 何时使用

以标签的形式展示用户选择或者输入的内容时

#### 基本使用

```jsx
const Button = require('../Button').default;
initialState =  { value: ['标签名字过长时省略显示内容','标签2','标签3'] };

<div>
    <ComboArea
        value={state.value}
        onChange={(value)=>setState({value})}
        placeholder={'请点击按钮添加一条新的标签'}
    />

    <div style={{ marginTop: '10px'}}>
        <Button
            width={200} 
            onClick={() => setState(() => ({
                value: [
                    ...state.value,
                    `标签${state.value.length + 1}`
                ]
            }))}>
            {'添加新标签'}
        </Button>
    </div>
</div>
```

#### 显示placeholder

```jsx
initialState =  { value: []};
<ComboArea
    value={state.value}
    placeholder={'placeholder'}
    onChange={(value)=>setState({value})}
/>
```

#### 禁用状态

```jsx
initialState =  { value: ['标签1','标签2','标签3'] };

<ComboArea
    value={state.value}
    disabled={true}
    onChange={(value)=>setState({value})}
/>
```

#### 错误状态

```jsx
initialState =  { value: ['标签1','标签2','标签3']  };
<ComboArea
    value={state.value}
    status={'error'}
    onChange={(value)=>setState({value})}
/>
```

#### 传入对象数组类型的value，配合formatter格式化输出的标签信息

```jsx
initialState =  { value: [{key:1,name:'标签1'},{key:2,name:'标签2'},{key:3,name:'标签3'}] };
<ComboArea
    value={state.value}
    formatter={(tag)=> tag.name}
    onChange={(value)=>setState({value})}
/>
```
