### 这个控件叫什么

开关选择器

### 何时使用

需要表示开关状态或两种状态之间的切换时，切换 Switch 会直接触发状态改变。

### 示例

#### 1. 可用开启状态

```jsx
initialState = { checked: true };
<Switch onChange={({detail}) => setState({checked: detail})} checked={state.checked}/>

```

#### 2. 可用关闭状态

```jsx
initialState = { checked: false };
<Switch onChange={({detail}) => setState({checked: detail})} checked={state.checked}/>

```

#### 3. 禁用开启状态

```jsx
initialState = { checked: true };
<Switch
    checked={state.checked}
    disabled={true} 
    onChange={({detail}) => setState({checked: detail})} 
/>

```

#### 4. 禁用关闭状态

```jsx
initialState = { checked: false };
<Switch
    checked={state.checked}
    disabled={true} 
    onChange={({detail}) => setState({checked: detail})} 
/>

