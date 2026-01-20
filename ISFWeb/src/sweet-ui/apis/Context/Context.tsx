import React from 'react'

interface ContextProviderProps extends React.ClassAttributes<void> {
    /**
     * Context将要传递的值
     */
    value?: any;
}

/**
 * ReactContext的简单实现
 * @param defaultValue 默认值
 * @example
 ```tsx
const { Provider, Consumer } = createContext('light')

class App extends React.Component {
    render() {
        return (
            <Provider
                value={'dark'}
            >
                <ContextExample />
            </Provider>
        )
    }
}

class ContextExample extends React.Component {
    render() {
        return (
            <Consumer>
                {
                    (theme) => <Button theme={theme}>Click Me</Button>
                }
            </Consumer>
        )
    }
}
 ```
 */
export function createContext(defaultValue?: any) {
    let context: any = defaultValue

    return {
        Provider: class Provider extends React.Component<ContextProviderProps, any> {
            componentDidMount() {
                const { value } = this.props

                context = value
            }

            componentDidUpdate(prevProps, prevState) {
                if (this.props.value !== this.context) {
                    context = this.props.value
                }
            }

            render() {
                return this.props.children
            }
        },

        Consumer: class Consumer extends React.Component {
            render() {
                const ContextedComponent = this.props.children

                if (typeof ContextedComponent === 'function') {
                    return ContextedComponent(context)
                } else {
                    return ContextedComponent
                }
            }
        },
    }
}