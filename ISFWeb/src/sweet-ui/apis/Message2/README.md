### 这个控件叫什么

弹窗提示

### 何时使用

常用于页面的错误反馈

### 示例

#### 1. 仅包含确定按钮弹窗（默认情况）

```jsx
    <div>
        <Button
            onClick={()=> Message2.info({message:'一般提示。'})}
        >
            {'提示窗'}
        </Button>
        <Button
            style ={{marginLeft:"20px"}}
            onClick={()=> Message2.alert({message:'一般提示。'})}
        >
            {'警告窗'}
        </Button>
        <Button
            style ={{marginLeft:"20px"}}
            onClick={()=> Message2.error({message:'一般提示。'})}
        >
            {'错误窗'}
        </Button>
    </div>
```

#### 2. 包含关闭按钮和确定/取消按钮的弹窗

* 关闭和取消按钮同时出现
* 当需要显示时，传入showCancelIcon为true

```jsx
const SweetIcon = require('../../components/SweetIcon').default;
    <div>
        <Button
            onClick={()=>  Message2.info({ message:'一般提示。',showCancelIcon:true})}
        >
            {'提示窗'}
        </Button>

        <Button
            style ={{marginLeft:"20px"}}
            onClick={()=> Message2.alert({ message:'一般提示。',showCancelIcon:true})}
        >
            {'警告窗'}
        </Button>

        <Button
            style ={{marginLeft:"20px"}}
            onClick={()=> Message2.error({ message:'一般提示。',showCancelIcon:true})}
        >
            {'错误窗'}
        </Button>
    </div>
```

#### 3. 包含红色文字提示的弹窗

* 需要显示红色文字提示时，传入title

```jsx
    <div>
        <Button
            onClick={()=> Message2.info({title:'无法执行xx操作',message:'原因。'})}
        >
            {'提示窗'}
        </Button>
        <Button
            style ={{marginLeft:"20px"}}
            onClick={()=> Message2.alert({title:'无法执行xx操作',message:'原因。'})}
        >
            {'警告窗'}
        </Button>
        <Button
            style ={{marginLeft:"20px"}}
            onClick={()=> Message2.error({title:'无法执行xx操作',message:'原因。'})}
        >
            {'错误窗'}
        </Button>
    </div>
```


#### 4. 包含错误详情的弹窗

* 需要显示错误详情时，传入detail

```jsx
    <div>
        <Button
            onClick={()=> Message2.info({ message:'错误信息', detail:"原因XXX"})}
        >
            {'提示窗'}
        </Button>
        <Button
            style ={{marginLeft:"20px"}}
            onClick={()=> Message2.alert({message:'错误信息', detail:<><div>{'111'}</div><div>{'222'}</div></>})}
        >
            {'警告窗'}
        </Button>
        <Button
            style ={{marginLeft:"20px"}}
            onClick={()=> Message2.error({message:'错误信息',detail:"原因XXX"})}
        >
            {'错误窗'}
        </Button>
    </div>
```

#### 5. 成功提示窗

```jsx
    <Button
        onClick={()=> Message2.success({ title:'成功提示',message:'原因或主注释。'})}
    >
    {'成功提示'}
    </Button>
```