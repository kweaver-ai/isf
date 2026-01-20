#### 何时使用

弹出一个下拉菜单给用户选择操作

#### 基本使用

```jsx
initialState = { value: 1 };
<Select 
    width={120}
    value={state.value}
    onChange={(event) => setState({value: event.detail})}
>
    <Select.Option value={1} >正常</Select.Option>
    <Select.Option value={2} >选中</Select.Option>
    <Select.Option value={3} disabled={true} >禁用</Select.Option>

</Select>

```

```jsx
<Select 
    width={120}
    disabled={true}
    placeholder={'禁用状态'}
>
    <Select.Option value={1} >正常</Select.Option>
    <Select.Option value={2} >选中</Select.Option>
    <Select.Option value={3} disabled={true} >禁用</Select.Option>

</Select>

```


