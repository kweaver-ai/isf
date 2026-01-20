#### 何时使用

#### 基本使用

```jsx
<Steps
    items={[
        {
            title: '第一步',
            description: '事实上事实上少时诵诗书飒飒飒',
        },
        {
            title: '第二步',
            description: '哈哈哈哈',
            subTitle: 'llll: 999:00',
        },
        {
            title: '第三步',
            status: 'error',
        }
    ]}
    current={1}
/>
```
```jsx
<Steps
    items={[
        {
            title: '第一步',
            description: '事实上事实上少时诵诗书飒飒飒',
        },
        {
            title: '第二步',
            description: '哈哈哈哈',
            subTitle: 'llll: 999:00',
        },
        {
            title: '第三步',
            status: 'error',
        }
    ]}
    direction={'vertical'}
    current={1}
/>
```


#### 可点击

```jsx
initialState = {
    current: 0,
};
<Steps
    items={[
        {
            title: '第一步',
        },
        {
            title: '第二步',
            subTitle: 'llll: 999:00',
        },
        {
            title: '第三步',
            status: 'error',
        },
        {
            title: '第四步',
            disabled: true,
        }
    ]}
    current={state.current}
    onChange={(v) => setState({current: v})}
/>
```

### small-size
```jsx
initialState = {
    current: 0,
};
<Steps
    items={[
        {
            title: '第一步',
        },
        {
            title: '第二步',
            subTitle: 'llll: 999:00',
        },
        {
            title: '第三步',
            status: 'error',
        },
        {
            title: '第四步',
            disabled: true,
        }
    ]}
    size={'small'}
    current={state.current}
    onChange={(v) => setState({current: v})}
/>
```

#### 纵向

```jsx
initialState = {
    current: 0,
};
<Steps
    items={[
        {
            title: '第一步',
            description: '事实上事实上少时诵诗书飒飒飒',
        },
        {
            title: '第二步',
            description: '哈哈哈哈',
            subTitle: 'llll: 999:00',
        },
        {
            title: '第三步',
        }
    ]}
    direction={'vertical'}
    current={state.current}
    onChange={(v) => setState({current: v})}
/>
```

### small-size
```jsx
initialState = {
    current: 0,
};
<Steps
    items={[
        {
            title: '第一步',
            description: '哈哈哈哈',
            subTitle: 'llll: 999:00',
        },
        {
            title: '第二步',
            subTitle: 'llll: 999:00',
        },
        {
            title: '第三步',
            status: 'error',
        },
        {
            title: '第四步',
            disabled: true,
        }
    ]}
    size={'small'}
    direction={'vertical'}
    current={state.current}
    onChange={(v) => setState({current: v})}
/>
```