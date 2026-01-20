### 这个控件叫什么

toast 提示

### 何时使用

常用于操作结果的反馈，提示语一般是短语

### 示例

#### 1. 基本使用 传入文本

```jsx
function showToast() {
    Toast.open('编辑成功')
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button>{'传入文本'}</Button>
</div>

```

#### 2. 传入jsx

```jsx
function showToast() {
    Toast.open(<span style={{ color: 'red' }}>{'新建成功'}</span>)
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button>{'传入jsx'}</Button>
</div>

```

#### 3. 传入duration为10000ms

```jsx
function showToast() {
    Toast.open('duration为10000ms', { duration: 10000 })
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button width={120}>{'传入duration'}</Button>
</div>

```

#### 4. 成功提示，深色

```jsx
function showToast() {
    Toast.success('深色成功提示')
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button width={120}>{'深色成功提示'}</Button>
</div>

```

#### 5. 错误提示，深色

```jsx
function showToast() {
    Toast.error('深色错误提示')
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button width={120}>{'深色错误提示'}</Button>
</div>

```

#### 6. 警告提示，深色

```jsx
function showToast() {
    Toast.warning('深色警告提示')
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button width={120}>{'深色警告提示'}</Button>
</div>

```

#### 7. 一般消息提示，深色

```jsx
function showToast() {
    Toast.info('深色消息提示')
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button width={120}>{'深色消息提示'}</Button>
</div>

```

#### 8. toast无图标，浅色

```jsx
function showToast() {
    Toast.lightOpen('白色背景')
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button width={120}>{'白色背景'}</Button>
</div>

```

#### 9. 成功提示，白色背景

```jsx
function showToast() {
    Toast.lightSuccess('浅色成功提示')
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button width={120}>{'浅色成功提示'}</Button>
</div>

```
#### 10. 错误提示，白色背景

```jsx
function showToast() {
    Toast.lightError('浅色错误提示')
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button width={120}>{'浅色错误提示'}</Button>
</div>

```

#### 11. 警告提示，白色背景

```jsx
function showToast() {
    Toast.lightWarning('浅色警告提示')
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button width={120}>{'浅色警告提示'}</Button>
</div>

```

#### 12. 一般消息提示，白色背景

```jsx
function showToast() {
    Toast.lightInfo('消息提示')
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={showToast}
>
    <Button width={120}>{'浅色消息提示'}</Button>
</div>

```

#### 13. 销毁

```jsx
function destoryToast() {
    Toast.destory()
}

<div
    style={{ display: 'inline-block', marginRight: '10px'}}
    onClick={destoryToast}
>
    <Button>{'销毁'}</Button>
</div>

```