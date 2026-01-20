### 这个控件叫什么

单选框

### 何时使用

用于在多个备选项中选中单个状态。

### 使用示例

#### 1.基本使用

```jsx
initialState = { checked: false };
<Radio checked={state.checked} onChange={({detail: {event, value}}) => setState({checked: event.target.checked})} >单选框</Radio>
```

```jsx
<Radio defaultChecked>初始状态为选中状态</Radio>

```

```jsx
<div>
<div style={{ display: 'inline-block', marginRight: '10px'}}>
<Radio disabled >禁用</Radio>
</div>
<Radio checked disabled >禁用</Radio>
</div>

```

#### 2.多个单选框

```jsx
initialState = { value: 1 };
<div>
<div style={{ display: 'inline-block', marginRight: '10px'}}>
<Radio name={'radio'} value={1} onChange={({detail: {event, value}}) => setState({value})} checked={state.value === 1}/>
</div>
<div style={{ display: 'inline-block', marginRight: '10px'}}>
<Radio name={'radio'} value={2} onChange={({detail: {event, value}}) => setState({value})} checked={state.value === 2}/>
</div>
<Radio name={'radio'} value={3} onChange={({detail: {event, value}}) => setState({value})} checked={state.value === 3}/>
</div>
```
