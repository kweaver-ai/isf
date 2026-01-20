### 这个控件叫什么

弹窗

### 何时使用

* 以弹窗的形式展示页面信息

* 标题，顶部的关闭按钮，内容区，底部的操作按钮均可根据实际使用情况传入

### 使用示例

#### 1.带有标题、顶部关闭按钮、确定取消按钮的复合弹窗

```jsx
const SweetIcon = require('../SweetIcon').default;
const Button = require('../Button').default;

initialState = { show: false };
<div>
    {
        state.show ?
            <ModalDialog2
                title={'内部共享'}
                width={760}
                icons={[{
                    icon: <SweetIcon name="x" size={16} />,
                    onClick: () => setState({show:false})
                    }
                ]}
                buttons={[ {
                    text: '确定',
                    theme: 'oem',
                    onClick: () =>setState({show:false}),
                },
                {
                    text: '取消',
                    theme: 'regular',
                    onClick: () =>setState({show:false}),
                } ]}
            >
                <div style={{ height: '500px'}}>
                    {'弹窗内容自定义'}
                </div>
            </ModalDialog2>
            :null
    }
    <Button onClick={() => setState({show:true})}>{'打开弹窗'}</Button>
</div>
```

#### 2.带有标题、顶部关闭按钮、底部关闭操作按钮的复合弹窗

```jsx
const SweetIcon = require('../SweetIcon').default;
const Button = require('../Button').default;

initialState = { show: false, validateResult: true };
<div>
    {
        state.show ?
            <ModalDialog2 
                title={'查看大小'}
                width={400}
                icons={[  {
                    icon: <SweetIcon name="x" size={16} />,
                    onClick: () => setState({show:false})
                    }
                ]}
                buttons={[ {
                    text: '关闭',
                    theme: 'oem',
                    onClick: () =>setState({show:false}),
                } ]}
            >
                <div>{'测试文件.docx'}</div>
                <div style={{ marginTop: '30px'}}>{'大小：1.23GB'}</div>
            </ModalDialog2>:null
    }
    <Button onClick={() => setState({show:true})}>{'打开弹窗'}</Button>
</div>
```

#### 3.仅有标题和关闭按钮的展示弹窗

```jsx
const SweetIcon = require('../SweetIcon').default;
const Button = require('../Button').default;

initialState = { show: false, validateResult: true };
<div>
    {
        state.show ?
            <ModalDialog2 
                title={'标题'}
                width={400}
                icons={[  {
                    icon: <SweetIcon name="x" size={16} />,
                    onClick: () => setState({show:false})
                    }
                ]}
            >
                <div style={{height:'112px'}}>这是一个展示信息这是一个展示信息这是一个展示信息。</div>
            </ModalDialog2>:null
    }
    <Button onClick={() => setState({show:true})}>{'打开弹窗'}</Button>
</div>
```

#### 4.标题居中和带有关闭按钮的展示弹窗

```jsx
const SweetIcon = require('../SweetIcon').default;
const Button = require('../Button').default;

initialState = { show: false, validateResult: true };
<div>
    {
        state.show ?
            <ModalDialog2 
                title={'标题（非必须）'}
                titleTextAlign={'center'}
                width={400}
                icons={[  {
                    icon: <SweetIcon name="x" size={16} />,
                    onClick: () => setState({show:false})
                    }
                ]}
            >
                <div style={{height:'120px'}}>这是一个展示信息这是一个展示信息这是一个展示信息。</div>
            </ModalDialog2>:null
    }
    <Button onClick={() => setState({show:true})}>{'打开弹窗'}</Button>
</div>
```

#### 5.不允许拖拽的弹窗

* 传入draggable为false，弹窗不允许拖拽

```jsx
const SweetIcon = require('../SweetIcon').default;
const Button = require('../Button').default;

initialState = { show: false, validateResult: true };
<div>
    {
        state.show ?
            <ModalDialog2 
                title={'标题'}
                draggable={false}
                width={400}
                icons={[  {
                    icon: <SweetIcon name="x" size={16} />,
                    onClick: () => setState({show:false})
                    }
                ]}
            >
                 <div style={{height:'112px'}}>这是一个展示信息这是一个展示信息这是一个展示信息。</div>
            </ModalDialog2>:null
    }
    <Button onClick={() => setState({show:true})}>{'打开弹窗'}</Button>
</div>
```