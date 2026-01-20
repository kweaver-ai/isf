### 这个控件叫什么

抽屉，屏幕边缘滑出的浮层面板。

### 何时使用

抽屉从父窗体边缘滑入，覆盖住部分父窗体内容。用户在抽屉内操作时不必离开当前任务，操作完成后，可以平滑地回到到原任务。

* 当需要一个附加的面板来控制父窗体内容，这个面板在需要时呼出；

* 当需要在当前任务流中插入临时任务，创建或预览附加内容。

### 使用示例

```jsx
const Button = require('../Button').default;

initialState = { open: false };
<div>
<Button onClick={() => setState({open: !state.open})}>{"Right"}</Button>
<Drawer 
    open={state.open}
    title={'Drawer Title'}
    footer={
        <div
            style={{
              position: 'absolute',
              right: 0,
              bottom: 0,
              padding: '13px 16px'
            }}
          >
            <Button onClick={() => setState({open: false})} theme="oem" style={{ marginRight: 8 }}>
              {'确定'}
            </Button>
            <Button onClick={() => setState({open: false})}>
              {'取消'}
            </Button></div>
    }
    onDrawerClose={(event) => setState({open: event.detail})}
>
    <div style={{textAlign: 'center'}}>{'Some contents'}</div>
</Drawer>
</div>
```

```jsx
const Button = require('../Button').default;

initialState = { open: false };
<div>
<Button onClick={() => setState({open: !state.open})}>{"Bottom"}</Button>
<Drawer 
    open={state.open}
    title={'Drawer Title'}
    position={'bottom'}
    destroyOnClose={true}
    footer={
        <div
            style={{
              position: 'absolute',
              right: 0,
              bottom: 0,
              padding: '13px 16px',
            }}
          >
            <Button onClick={() => setState({open: false})} theme="oem" style={{ marginRight: 8 }}>
              {'确定'}
            </Button>
            <Button onClick={() => setState({open: false})}>
              {'取消'}
            </Button></div>
    }
    onDrawerClose={(event) => setState({open: event.detail})}
>
    <div style={{textAlign: 'center'}}>{'Some contents'}</div>
</Drawer>
</div>
```

```jsx
const Button = require('../Button').default;

initialState = { open: false };
<div>
<Button onClick={() => setState({open: !state.open})}>{"Left"}</Button>
<Drawer 
    open={state.open}
    title={'Drawer Title'}
    position={'left'}
    showMask={false}
    footer={
        <div
            style={{
              position: 'absolute',
              right: 0,
              bottom: 0,
              padding: '13px 16px',
            }}
          >
            <Button onClick={() => setState({open: false})} theme="oem" style={{ marginRight: 8 }}>
              {'确定'}
            </Button>
            <Button onClick={() => setState({open: false})}>
              {'取消'}
            </Button></div>
    }
    onDrawerClose={(event) => setState({open: event.detail})}
>
    <div style={{textAlign: 'center'}}>{'Some contents'}</div>
</Drawer>
</div>
```

```jsx
const Button = require('../Button').default;

initialState = { open: false };
<div>
<Button onClick={() => setState({open: !state.open})}>{"Top"}</Button>
<Drawer 
    open={state.open}
    title={'Drawer Title'}
    position={'top'}
    footer={
        <div
            style={{
              position: 'absolute',
              right: 0,
              bottom: 0,
              padding: '13px 16px'
            }}
          >
            <Button onClick={() => setState({open: false})} theme="oem" style={{ marginRight: 8 }}>
              {'确定'}
            </Button>
            <Button onClick={() => setState({open: false})}>
              {'取消'}
            </Button></div>
    }
    onDrawerClose={(event) => setState({open: event.detail})}
>
    <div style={{textAlign: 'center'}}>{'Some contents'}</div>
</Drawer>
</div>
```