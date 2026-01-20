### 何时使用

封装了一个操作命令，响应用户的点击行为，触发响应的业务逻辑。

### 示例

#### 按钮尺寸

* 按钮尺寸默认为80x30px，固定高度30px，宽度可调整（为8的倍数）
* 设置`size:auto`，按钮宽度按内容自适应
* 设置`width, height`指定具体宽高

```jsx
<div>
    <div style={{ display: 'inline-block', marginRight: '10px'}}>
        <Button>{'常规'}</Button>
    </div>
    <div style={{ display: 'inline-block',marginRight: '10px'}}>
        <Button size={'auto'}>{'适应字体长短'}</Button>
    </div>
    <Button  width={360} height={40}>{'指定宽高'}</Button>
</div>

```


#### 1. 一般按钮

##### 白色按钮

```jsx
<div>
    <div style={{ display: 'inline-block', marginRight: '10px'}}>
    <Button onClick={(event) => alert('you clicked a regular button.')}>{'正常'}</Button>
    </div>
    <Button  disabled={true} >{'禁用'}</Button>
</div>

```

##### 主题色（可OEM）按钮

```jsx
<div>
    <div style={{ display: 'inline-block', marginRight: '10px'}}>
        <Button theme={'oem'} onClick={(event) => alert('you clicked an oem-style button.')}>{'正常'}</Button>
    </div>
    <Button theme={'oem'} disabled={true} >{'禁用'}</Button>
</div>
```

##### 灰色按钮

```jsx
<div>
    <div style={{ display: 'inline-block', marginRight: '10px'}}>
        <Button theme={'gray'} onClick={(event) => alert('you clicked a gray button.')}>{'正常'}</Button>
    </div>
    <Button theme={'gray'} disabled={true} >{'禁用'}</Button>
</div>
```

##### 深色背景按钮

```jsx
<div style={{ width: '260px', height: '32px', padding: '10px', borderRadius: '4px', backgroundColor: '#333' }}>
    <div style={{ display: 'inline-block', marginRight: '10px'}}>
    <Button theme={'dark'} onClick={(event) => alert('you clicked a dark button.')}>{'正常'}</Button>
    </div>
    <Button theme={'dark'} disabled={true} >{'禁用'}</Button>
</div>
```

#### 2. 图标按钮

* 当需要在`Button`内嵌入图标时，可以设置`icon`属性指定图标名称或图标元素

```jsx
<div>
    <div style={{ display: 'inline-block', marginRight: '10px'}}>
       <Button  icon={'search'}>{'按钮'}</Button>
    </div>
    <Button icon={'search'} disabled={true}>{'禁用'}</Button>
</div>

```

#### 3. 文字按钮

```jsx
<div>
    <div style={{ display: 'inline-block', marginRight: '20px'}}>
    <Button theme={'text'} onClick={(event) => alert('you clicked a text button.')}>{'正常'}</Button>
    </div>
    <Button theme={'text'} disabled={true} >{'禁用'}</Button>
</div>
```

#### 4. 无边框按钮

多用于工具栏按钮

```jsx
<div>
    <div style={{ display: 'inline-block', marginRight: '20px'}}>
        <Button theme={'noborder'} icon={'search'}>{'工具栏按钮'}</Button>
    </div>
    <Button theme={'noborder'} icon={'search'} disabled={true} >{'禁用'}</Button>
</div>

```

