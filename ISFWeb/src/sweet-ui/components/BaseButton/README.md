# BaseButton

基本按钮控件，无样式。

## Usage

```jsx
<BaseButton
    disabled={true}
    style={{
        color: '#000',
        fontSize: '20px',
    }}
    className="customizeClassName"
    onClick={event => {}}
>
    {'BaseButton'}
</BaseButton>
```

## Props

| 参数                | 说明                       | 必须 | 类型                | 默认值 |
| ------------------- | -------------------------- | ---- | ------------------- | ------ |
| disabled            | 控制按钮禁用               | No   | boolean             | false  |
| style               | 自定义 style 样式          | No   | React.CSSProperties | /      |
| className           | 自定义组件根节点 className | No   | string              | /      |
| [onClick](#onClick) | 按钮点击时的处理函数       | No   | function            | /      |

## Reference

### Props

#### <span id="onClick">onClick</span>

按钮点击时调用的处理函数。

| 参数  | 说明         | 类型                                      |
| :---- | :----------- | :---------------------------------------- |
| event | 鼠标事件对象 | React.MouseEvent&lt;HTMLButtonElement&gt; |
