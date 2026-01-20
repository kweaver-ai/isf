### 这个控件叫什么

数字输入框

### 何时使用

* 显示纯数值时；

* 当需要输入一定范围的标准数值时。

### 基本使用

```jsx
<NumberBox
    placeholder={'请输入数字'}
    onValueChange={(event) => {setState({value: event.detail})} }
/>
```

#### 1. 只允许输入正数

* 设置 `min={0}`，只能输入正数，不允许输入负数

```jsx
initialState = { value: 8 };
<NumberBox
    value={state.value}
    min={0}
    onValueChange={(event) => setState({value: event.detail})}
/>
```

#### 2. 只允许输入整数

* 设置精度 `precision={0}`，只能输入整数，不允许输入小数

```jsx
initialState = { value: 3 };
<NumberBox
    value={state.value}
    precision={0}
    onValueChange={(event) => setState({value: event.detail})}
/>
```

#### 3. 允许输入小数

* 默认设置中，允许输入小数，且无精度限制，可输入任意位数的小数

```jsx
initialState = { value: 1.2 };

<NumberBox
    value={state.value} 
    onValueChange={(event) => setState({value: event.detail})}
/>
```

* 可通过`precision`属性指定小数值精度

* `precision={1}`数字框显示为1位小数

```jsx
initialState = { value: 1.0 };

<NumberBox
    precision={1}
    value={state.value} 
    onValueChange={(event) => setState({value: event.detail})}
/>
```

* `precision={3}`数字框显示为3位小数

```jsx
initialState = { value: 1.200 };

<NumberBox
    precision={3}
    value={state.value} 
    onValueChange={(event) => setState({value: event.detail})}
/>
```

#### 4. 允许使用步进（鼠标点击上下箭头或者键盘键入上下箭头改变输入框的值）

* 可通过`step`属性指定每次增量的大小

* `step={0}`，不显示上下键按钮且无法通过键盘改变输入框的值

```jsx
initialState = { value: 1.2 };

<NumberBox
    step={0}
    value={state.value} 
    onValueChange={(event) => setState({value: event.detail})}
/>
```

* `step={0.1}`，点击上下键每次增加或者减小的值为0.1

```jsx
initialState = { value: 1.2 };

<NumberBox
    step={0.1}
    value={state.value} 
    onValueChange={(event) =>{
        setState({value: event.detail})
    } }
/>
```

* `step={1}`，点击上下键每次增加或者减小的值为1

```jsx
initialState = { value: 1.2 };

<NumberBox
    step={1}
    value={state.value} 
    onValueChange={(event) => setState({value: event.detail})}
/>
```

* `step={1}`，点击上下键每次增加或者减小的值为1
* 传入最大最小值，当输入框的内容为不完整输入（如空），点击向上从最小值开始递增，点击向下，从最大值开始递减

```jsx
initialState = { value: 1 };

<NumberBox
    step={1}
    min={0}
    max={20}
    value={state.value} 
    onValueChange={(event) => setState({value: event.detail})}
/>
```

* `step={5}`，点击上下键每次增加或者减小的值为5

```jsx
initialState = { value: 1.2 };

<NumberBox
    step={5}
    value={state.value} 
    onValueChange={(event) => setState({value: event.detail})}
/>
```

* 当同时设置步进和精度时候，如果步进（0.01）小于精度限制（1位小数），则步进不生效

#### 5. 指定最大值最小值

* 默认限制最大值不大于Number.MAX_SAFE_INTEGER，最小值不小于Number.MIN_SAFE_INTEGER

```jsx
initialState = { value:2 };
<NumberBox
    value={state.value} 
    onValueChange={(event) => setState({value: event.detail})}
/>
```

* 当设置最大值为10时候，允许输入超过10的数字，但是失焦时，还原成最大值10

```jsx
initialState = { value: 4 };
<NumberBox
    max={10}
    value={state.value} 
    onValueChange={(event) => setState({value: event.detail})}
/>
```

* 当设置最小值设置为2时，允许输入小于2的数字（无法键入负数），但是失焦时，还原成最小值2

```jsx
initialState = { value: 4 };
<NumberBox
    min={2}
    value={state.value} 
    onValueChange={(event) => setState({value: event.detail})}
/>
```

#### 6. 禁用

```jsx
<NumberBox disabled={true} value={10}/>

```

#### 7. 只读

```jsx
<NumberBox readOnly={true} value={10}/>

```
