### 这个控件叫什么

标签

### 何时使用

标记事物的属性或者事物分类时使用

### 示例

#### 1. 基本用法

```jsx

<Tag>{'标签'}</Tag>

```

#### 2. 可关闭的标签

* 传入 closable 为 true 标签显示关闭按钮，支持关闭操作，点击关闭按钮时，向外抛出 onClose 事件

```jsx
initialState = { tags: [{key:1,name:'标签1'},{key:2,name:'标签2'},{key:3,name:'标签3'}] };

function handleClose(key){
    setState({
        tags:state.tags.filter((tag)=>{
            return tag.key !== key
        })
    })
}

<div style={{height: '30px',lineHeight:'30px'}}>
    {
        state.tags.map((tag)=>{
            return <Tag style={{marginRight:'8px'}} key={tag.key} closable={true} onClose={()=>handleClose(tag.key)}>{tag.name}</Tag>
        })
    }
</div>
```

#### 3. 可选择的标签

* 传入 checkable 为 true 标签可选择，支持选择操作，点击标签时，向外抛出 onChange 事件

```jsx
initialState = {  tags: [{key:1,checked:false,name:'标签1'},{key:2,checked:false,name:'标签2'},{key:3,checked:false,name:'标签3'}] };

function handleChange(key){
    setState({
        tags:state.tags.map((tag)=>{
            return key === tag.key ? {...tag,checked:!tag.checked,}:tag
        })
    })
}

<div style={{height: '30px',lineHeight:'30px'}}>
    {
        state.tags.map((tag)=>{
            return <Tag 
                style={{marginRight:'8px'}}
                key={tag.key} 
                checkable={true}
                checked={tag.checked}
                onChange={()=>handleChange(tag.key)}>{tag.name}
            </Tag>
        })
    }
</div>

```

#### 4. 带删除按钮的可选择的标签

* 点击删除按钮时，触发 onClose 关闭事件，点击其他区域触发 onChange 选择事件

```jsx
initialState = {  tags: [{key:1,checked:false,name:'标签1'},{key:2,checked:false,name:'标签2'},{key:3,checked:false,name:'标签3'}] };

function handleClose(key){
    alert('目前点击的是关闭按钮的操作区，关闭弹窗后，该标签被删除')
    setState({
        tags:state.tags.filter((tag)=>{
            return tag.key !== key
        })
    })
}

function handleChange(key){
    alert('目前点击的是非关闭按钮的操作区，关闭弹窗后，该标签的选中状态发生改变')
    setState({
        tags:state.tags.map((tag)=>{
            return key === tag.key ? {...tag,checked:!tag.checked,}:tag
        })
    })
}

<div style={{height: '30px',lineHeight:'30px'}}>
    {
        state.tags.map((tag)=>{
            return <Tag 
                style={{marginRight:'8px'}}
                key={tag.key} 
                checkable={true}
                closable={true}
                checked={tag.checked}
                onClose={()=>handleClose(tag.key)}
                onChange={()=>handleChange(tag.key)}
            >
                {tag.name}
            </Tag>
        })
    }
</div>
```

#### 5. 禁用

* 传入 disabled 为 true，标签不支持选择事件和关闭事件

```jsx
initialState = {  checked:false };

function handleChange(){
    alert('禁用状态，此弹窗无机会弹出')
}

function handleClose(){
    alert('禁用状态，此弹窗无机会弹出')
}

<Tag 
    disabled={true} 
    checkable={true} 
    closable={true}
    checked={state.checked}
    onChange={()=>handleChange()}
    onClose={()=>handleClose()}
>
    {'标签'}
</Tag>
```

#### 6. 大尺寸标签

* 通过传入 style 或者 className 改变标签大小

```jsx
initialState = { checked:false };

function handleChange(){
    setState({
        checked:!state.checked
    })
}

<Tag 
    style={{height:'30px',lineHeight:'30px',paddingLeft:'16px',paddingRight:'16px'}} 
    checkable={true}
    checked={state.checked}
    onChange={()=>handleChange()}
>
    {'大标签'}
</Tag>
```
