### 这个控件叫什么

下拉按钮

### 何时使用

配合下拉菜单一起使用

### 示例

#### 1.下拉白色

```jsx

<div>
    <div style={{ display: 'inline-block', marginRight: '10px'}}>
        <DropButton>{'正常'}</DropButton>
    </div>
    <DropButton disabled={true}>{'禁用'}</DropButton>
</div>

```

#### 2.下拉灰色

设置`theme: gray`显示灰色按钮

```jsx

<div>
    <div style={{ display: 'inline-block', marginRight: '10px'}}>
        <DropButton theme={'gray'}>{'正常'}</DropButton>
    </div>
    <DropButton theme={'gray'} disabled={true}>{'禁用'}</DropButton>
</div>

```

#### 3.下拉文字按钮

* 默认操作对象悬停有背景

```jsx

<div>
    <div style={{ display: 'inline-block', marginRight: '10px'}}>
        <DropButton theme={'text'} >{'正常'}</DropButton>
    </div>
    <DropButton theme={'text'} disabled={true}>{'禁用'}</DropButton>
</div>

```

* 设置`background:false`，操作对象悬停无背景

```jsx

<div>
    <div style={{ display: 'inline-block', marginRight: '10px'}}>
        <DropButton theme={'text'} background={false}>{'正常'}</DropButton>
    </div>
    <DropButton theme={'text'} background={false} disabled={true}>{'禁用'}</DropButton>
</div>

```
