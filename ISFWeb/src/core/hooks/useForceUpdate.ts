import * as React from 'react'

const { useState, useCallback } = React

function useForceUpdate() {
    const [, setValue] = useState<number>(0)

    return useCallback(() => {
        setValue((value) => value + 1)
    }, [])
}

export {
    useForceUpdate,
}