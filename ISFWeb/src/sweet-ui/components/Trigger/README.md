#### 点击触发

```jsx
const Button = require('../Button').default;
const Menu = require('../Menu').default;

<Trigger
    triggerEvent={'click'}
    anchorOrigin={[ 'left', 'bottom' ]}
    alignOrigin={[ 'left', 'top' ]}
    renderer={({setPopupVisibleOnClick}) => 
        <Button 
            onClick={setPopupVisibleOnClick}
            style={{width: '100px', height: '32px'}} 
        >{'Click me'}</Button>
    }
    freeze={false}
    onBeforePopupClose={(event) => 
        {
            // 阻止默认关闭事件
            //event.preventDefault()
        }} 
>
    {
        ({close}) => <Menu width={120}>
            <Menu.Item value={1} onClick={close}>正常</Menu.Item>
            <Menu.Item value={2} selected={true} onClick={close}>选中</Menu.Item>
            <Menu.Item value={3} disabled={true} >禁用</Menu.Item>
        </Menu>      
    } 
</Trigger>
```

#### 鼠标悬浮触发

```jsx
const Button = require('../Button').default;
const Menu = require('../Menu').default;
<Trigger
    triggerEvent={'hover'}
    anchorOrigin={[ 'left', 'bottom' ]}
    alignOrigin={[ 'left', 'top' ]}
    freeze={false}
    renderer={({setPopupVisibleOnMouseEnter, setPopupVisibleOnMouseLeave}) =>  
        <Button 
            onMouseEnter={setPopupVisibleOnMouseEnter}
            onMouseLeave={setPopupVisibleOnMouseLeave}
            width={120}
        >
            Hover me
        </Button>
    }
>
    {
        ({close}) => <Menu width={120}>
            <Menu.Item value={1} onClick={close}>正常</Menu.Item>
            <Menu.Item value={2} selected={true} onClick={close}>选中</Menu.Item>
            <Menu.Item value={3} disabled={true} >禁用</Menu.Item>
        </Menu>      
    }
</Trigger>
```

#### 元素聚焦时触发

```jsx
const TextBox = require('../TextBox').default;
const Menu = require('../Menu').default;
<Trigger
    triggerEvent={'focus'}
    anchorOrigin={[ 'left', 'bottom' ]}
    alignOrigin={[ 'left', 'top' ]}
    renderer={({setPopupVisibleOnFocus, setPopupVisibleOnBlur}) =>  
        <TextBox onFocus={setPopupVisibleOnFocus} onBlur={setPopupVisibleOnBlur} placeholder={'focus me'}/>
    }
    freeze={false}
>
    {
        ({close}) =><Menu width={200}>
            <Menu.Item value={1} onClick={close}>正常</Menu.Item>
            <Menu.Item value={2} selected={true} onClick={close}>选中</Menu.Item>
            <Menu.Item value={3} disabled={true} >禁用</Menu.Item>
        </Menu>     
    }
</Trigger>
```

