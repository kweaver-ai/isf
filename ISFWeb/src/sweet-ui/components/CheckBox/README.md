### 这个组件叫什么

复选框
  
### 何时使用

* 在一组可选项中进行多项选择时

* 单独使用可以表示两种状态之间的切换，用于状态标记

### 使用示例

#### 1.基本使用

```jsx
initialState = { checked: false};
<div>
<div style={{ display: 'inline-block', marginRight: '30px'}}>
  <CheckBox
    checked={state.checked}
    className="CustomizeClassName"
    onCheckedChange={event => setState({checked: event.detail})}
/>
</div>
<div style={{ display: 'inline-block', marginRight: '30px'}}>
<CheckBox
    checked={true}
    className="CustomizeClassName"
/>
</div>
<CheckBox
    disabled
    onCheckedChange={event => setState({checked: event.detail})}
  />
</div>
```

#### 2.复选框+文字
```jsx
initialState = { checked: true};
<div>
<div style={{ display: 'inline-block', marginRight: '10px'}}>
<CheckBox
    checked={state.checked}
    onClick={event=>{}}
    onCheckedChange={event => setState({checked: event.detail})}
>
  {'复选框'}
</CheckBox>
</div>
<CheckBox
    checked={state.checked}
    disabled
    onCheckedChange={event => setState({checked: event.detail})}
>
  {'复选框'}
</CheckBox>
</div>
```

#### 3.半选
```jsx
<div>
<div style={{ display: 'inline-block', marginRight: '10px'}}>
<CheckBox
    indeterminate={true}
>
  {'复选框'}
</CheckBox>
</div>
<CheckBox
    indeterminate={true}
    disabled
>
  {'复选框'}
</CheckBox>
</div>
```