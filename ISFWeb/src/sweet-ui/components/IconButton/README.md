### 这个控件叫什么

纯图标按钮

### 何时使用

常用于列表中的可操作项

### 示例

#### 1. 基本使用

```jsx
const SweetIcon = require('../SweetIcon').default;

<IconButton
    icon={<SweetIcon name={'first'} title={'first'}/>}
    onClick={(event) => alert('You just clicked IconButton.', event)}
/>

```

#### 2. 禁用

```jsx
const SweetIcon = require('../SweetIcon').default;
<IconButton
    icon={<SweetIcon name={'first'}/>}
    disabled={true}
/>

```