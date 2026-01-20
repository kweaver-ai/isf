### 何时使用

* 需要用户输入表单域内容时

* 提供组合型输入框，带搜索的输入框等

### 基本使用

```jsx
initialState = { value: 'aaa' };
<TextBox 
    value={state.value} 
    onValueChange={(event) => {console.log(event.detail); setState({value: event.detail})} }
    placeholder={'请输入内容'}
    onBlur={() => {console.log('blur')}}
/>

```

```jsx
<TextBox 
    defaultValue={'bbb'} 
    onValueChange={(event) => console.log(event.detail)} 
    placeholder={'请输入内容'}
    onBlur={() => {console.log('blur')}}
/>

```

#### 禁用

```jsx
<TextBox disabled={true} placeholder={'请输入内容'}/>

```

#### 禁用状态下改变字体颜色

```jsx
<TextBox disabled={true} value={'禁用状态的红色文字显示'} placeholder={'请输入内容'}  style={{color:'red'}}/>

```

#### 指定宽度

```jsx
<TextBox width={300} placeholder={'请输入内容'}/>

```

#### 自动选中

> 指定`autoFocus:true`实现自动选中效果，指定`selectOnFocus:true`实现聚焦时选中输入内容的效果。

```jsx
initialState = { value: 'text' };
<TextBox 
    value={state.value} 
    autoFocus={true}
    selectOnFocus={true} 
    onValueChange={(event) => setState({value: event.detail})} 
    placeholder={'请输入内容'}
    onBlur={() => {console.log('blur')}}
/>

```


