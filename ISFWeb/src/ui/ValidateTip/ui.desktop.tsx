import React from 'react';
import Tip from '../Tip/ui.desktop';

export default ({ align = 'right', children }: { align?: string; children?: any }) => (
    <Tip borderColor="#CC9933" backgroundColor="#FFFFCC" align={align}>
        {
            children
        }
    </Tip>
)