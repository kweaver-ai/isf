### 这个控件叫什么

文本域

### 何时使用

用于多行输入

### 使用示例

#### 1.基本使用

```jsx
<TextArea 
    placeholder={'input here...'}
/>

```
```jsx
<TextArea 
    disabled={true}
    placeholder={'disabled'}
/>

```

```jsx
<TextArea 
    readOnly={true}
    value={'只读，不支持输入'}
/>

```

```jsx
initialState = { value: 'aaa' };
<TextArea 
    value={state.value} 
    onValueChange={(event) => {setState({value: event.detail})} }
    placeholder={'请输入内容'}
    onPressEnter={console.log}
/>

```

#### 2.自定义宽高

```jsx
<TextArea 
    placeholder={'input here...'}
    width={400}
    height={64}
/>

```

#### 3.限制输入字数

```jsx
<TextArea 
    placeholder={'input here...'}
    maxLength={500}
/>

```



#### 4.非受控，通过defaultValue指定初始值

```jsx
<TextArea 
    placeholder={'input here...'}
    defaultValue={'输入内容'}
/>

```
