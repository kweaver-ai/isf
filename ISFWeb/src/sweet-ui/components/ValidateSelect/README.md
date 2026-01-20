#### 何时使用

弹出一个下拉菜单给用户选择操作，发生错误时，有错误验证

#### 基本使用

```jsx
initialState = { value:1,data: [{value:1,label:'选项一'},{value:2,label:'选项二'},{value:3,label:'选项三'},{value:4,label:'无法通过验证的选项'}],validateState:'normal' };

<div>
    <ValidateSelect
        value={state.value}
        onChange={(event) => {
            setState({value: event.detail,validateState:'normal'})
        }}
        validateState={state.validateState}
        validateMessages={{
            ['empty']:'不允许选择超长项目',
        }}
    >
    {
        state.data.map(option=>{
            return <ValidateSelect.Option key={option.value} value={option.value} >{option.label}</ValidateSelect.Option>
        })
    }
    </ValidateSelect>

    <div style={{marginTop:'10px'}}>
        <Button 
            onClick={() => {
                if(state.value===4){
                    setState({validateState:'empty'})
                }else{
                    setState({validateState:'normal'})
                    alert('通过验证，你选择的内容将被提交')
                }
            }}
        >
            {'确定'}
        </Button>
    </div>
</div>
```
