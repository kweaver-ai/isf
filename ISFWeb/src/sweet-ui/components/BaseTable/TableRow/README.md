# TableRow

Table 的子组件，用于渲染表格中的一行，不带任何样式。

```react
<TableRow
    className="customizeClassName"
    hoverClassName="customizeHoverClassName"
    onClick={(event)=>{}}
    onEnter={(event)=>{}}
    onLeave={(event)=>{}}
>
    <td>cell1</td>
    <td>cell2</td>
    <td>cell3</td>
</TableRow>
```

## Props

| 参数                | 说明                       | 必须 | 类型     | 默认值 |
| ------------------- | -------------------------- | ---- | -------- | ------ |
| className           | 根节点 tr 的 className     | No   | string   | /      |
| hoverClassName      | 鼠标悬浮时的行的 className | No   | string   | /      |
| [onClick](#onClick) | 鼠标点击行时触发的处理函数 | No   | function | /      |
| [onEnter](#onEnter) | 鼠标移入行时触发的处理函数 | No   | function | /      |
| [onLeave](#onLeave) | 鼠标移出行时触发的处理函数 | no   | function | /      |

## Reference

### <span id="onClick">onClick</span>

鼠标点击行时触发的处理函数。

| 参数  |     说明     | 类型                                      |
| :---: | :----------: | ----------------------------------------- |
| event | 鼠标事件对象 | React.MouseEvent&lt;HTMLButtonElement&gt; |

### <span id="onEnter">onEnter</span>

鼠标移入行时触发的处理函数

| 参数  | 说明         | 类型                                      |
| ----- | ------------ | ----------------------------------------- |
| event | 鼠标事件对象 | React.MouseEvent&lt;HTMLButtonElement&gt; |

### <span id="onLeave">onLeave</span>

鼠标移出行时触发的处理函数

| 参数  | 说明         | 类型                                      |
| ----- | ------------ | ----------------------------------------- |
| event | 鼠标事件对象 | React.MouseEvent&lt;HTMLButtonElement&gt; |
