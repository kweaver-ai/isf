# FlatList

基础的列表组件，默认根据视口大小进行数据分片。

## Usage

```react
<FlatList
     height={300}
     getItemLayout={(data, index) => ({ length: 100, offset: 100 * (index + 1), index })}
     renderItem={(data, index) => (<div style={{ height: 100 }}>{`row${index} ${data}`}</div>)}
     refreshing={this.state.refreshing}
     RefreshingIndicatorComponent={<Button>Refreshing</Button>} // React实例
     ListEmptyComponent={()=><div>Empty</div>} // 无状态组件
     onEndReached={()=>{//do something}}
     data={[]}
/>
```

## Props

| 属性                                                          | 说明                                                 | 必须 | 类型                                              | 默认值 |
| ------------------------------------------------------------- | ---------------------------------------------------- | ---- | ------------------------------------------------- | ------ |
| data                                                          | 列表数据源                                           | Yes  | array                                             | /      |
| [renderItem](#renderItem)                                     | 定义如何渲染每一行数据                               | Yes  | function                                          | /      |
| [getItemLayout](#getItemLayout)                               | 获取每行数据的布局信息，用于计算最终要渲染的数据切片 | Yes  | function                                          | /      |
| height                                                        | 列表容器高度                                         | No   | number&#124;string                                | /      |
| [keyExtractor](#keyExtractor)                                 | 生成列表项的 key 值                                  | No   | function                                          | /      |
| [onScroll](#onScroll)                                         | 滚动列表时触发的处理函数                             | No   | function                                          | /      |
| [onEndReached](#onEndReached)                                 | 当列表滚动接近底部时触发的处理函数                   | No   | function                                          | /      |
| [ListEmptyComponent](#ListEmptyComponent)                     | 列表为空时显示                                       | No   | React.FunctionComponent &#124; React.ReactElement | /      |
| [RefreshingIndicatorComponent](#RefreshingIndicatorComponent) | 列表正在刷新时显示的指示信息                         | No   | React.FunctionComponent &#124; React.ReactElement | /      |
| refreshing                                                    | 指示列表是否正在刷新                                 | No   | boolean                                           | /      |

## Reference

### Props

#### <span id="renderItem">renderItem</span>

通过每一行的数据 record 以及当前行的索引 index 自定义数据渲染，例如根据不同数据和索引进行个性化显示。

| 参数  | 说明         | 类型   |
| ----- | ------------ | ------ |
| data  | 当前行的数据 | any    |
| index | 当前行的索引 | number |

#### <span id="getItemLayout">getItemLayout</span>

返回每一行的布局信息对象 ItemLayout。

| 参数  | 说明         | 类型   |
| ----- | ------------ | ------ |
| data  | 当前行的数据 | any    |
| index | 当前行的索引 | number |

布局对象 ItemLayout：

| 键     | 说明                             | 值类型 |
| ------ | -------------------------------- | ------ |
| length | 数据渲染的高度                   | number |
| offset | 该条数据渲染距离列表顶端的偏移量 | number |
| index  | 该条数据的下标                   | number |

#### <span id="keyExtractor">keyExtractor</span>

生成列表每一行的 key

| 参数  | 说明         | 类型   |
| ----- | ------------ | ------ |
| data  | 当前行的数据 | any    |
| index | 当前行的索引 | number |

#### <span id="onScroll">onScroll</span>

在列表容器内滚动时触发的处理函数。

| 参数  | 说明           | 类型                           |
| ----- | -------------- | ------------------------------ |
| event | 接收到事件对象 | React.UIEvent< HTMLDivElement> |

#### <span id="onEndReached">onEndReached</span>

当列表滚动接近底部时触发的处理函数。

#### <span id="ListEmptyComponent">ListEmptyComponent</span>

列表为空时显示的提示信息。可传入 React 实例或者无状态组件。

> 使用绝对定位铺满整个容器，样式需要外部定义。

#### <span id="RefreshingIndicatorComponent">RefreshingIndicatorComponent</span>

列表正在刷新时的指示信息。可传入 React 实例或者无状态组件。

> 使用绝对定位铺满整个容器，样式需要外部定义。
