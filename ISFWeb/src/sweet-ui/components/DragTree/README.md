#### 何时使用

树支持拖拽

#### 基本使用

```jsx
const SweetIcon = require('../SweetIcon').default;

const treeData = [
    {
        key: '1',
        title: '未分配组',
        isLeaf: true,
        icon: <SweetIcon name={'caution'}/>,
        data: {
            id: '1',
            isOrg: true,
            name: '未分配组',
            subDepCount: 0
        }
    },
    {
        key: '2',
        title: '全部用户',
        isLeaf: true,
        icon: <SweetIcon name={'caution'}/>,
        data: {
            id: '2',
            isOrg: true,
            name: '全部用户',
            subDepCount: 0
        }
    },
    {
        key: '3',
        title: '组织结构',
        isLeaf: true,
        icon: <SweetIcon name={'caution'}/>,
        data: {
            id: '3',
            isOrg: true,
            name: '组织结构',
            subDepCount: 0,
        }
    },
    {
        key: '4',
        title: 'hh奥斯迪阿萨帝大家爱到极啊啊啊啊啊啊啊啊啊',
        isLeaf: false,
        icon: <SweetIcon name={'caution'}/>,
        children: [
            {
                key: '41',
                title: 'bm',
                isLeaf: true,
                icon: <SweetIcon name={'edit'}/>,
                data: {
                    id: '41',
                    isOrg: false,
                    name: 'bm',
                    subDepCount: 0
                }
            },
            {
                key: '42',
                title: 'aaaddd',
                isLeaf: true,
                icon: <SweetIcon name={'edit'}/>,
                data: {
                    id: '42',
                    isOrg: false,
                    name: 'aaaddd',
                    subDepCount: 0
                }
            }
        ],
        data: {
            id: '4',
            isOrg: true,
            name: 'hh奥斯迪阿萨帝大家爱到极啊啊啊啊啊啊啊啊啊',
            subDepCount: 2
        }
    },
    {
        key: '5',
        title: 'wewaaaa',
        isLeaf: false,
        icon: <SweetIcon name={'caution'}/>,
        data: {
            id: '5',
            isOrg: true,
            name: 'wewaaaa',
            subDepCount: 1
        },
        children: [
            {
                key: '51',
                title: 'fffsd',
                isLeaf: true,
                icon: <SweetIcon name={'edit'}/>,
                data: {
                    id: '51',
                    isOrg: false,
                    name: 'fffsd',
                    subDepCount: 0
                },
            }
        ]
    },
    {
        key: '6',
        title: 'hhh',
        isLeaf: true,
        icon: <SweetIcon name={'caution'}/>,
        data: {
            id: '6',
            isOrg: true,
            name: 'hhh',
            subDepCount: 0,
        }
    },
    {
        key: '7',
        title: 'tttt',
        isLeaf: true,
        icon: <SweetIcon name={'caution'}/>,
        data: {
            id: '7',
            isOrg: true,
            name: 'tttt',
            subDepCount: 0,
        }
    },
];

<div style={{width: 270, height: 500, overflow: 'auto'}}>
    <DragTree treeData={treeData}/>
</div>
```