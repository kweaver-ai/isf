import React, { useState, useContext } from 'react';
import { useDrag, useDrop } from 'react-dnd';
import DragIcon from '../../icons/drag.svg';
import AppConfigContext from '@/core/context/AppConfigContext';

const ItemTypes = {
  CARD: 'card',
};

interface DraggableCardProps {
  id: string | number;
  index: number;
  moveCard: (fromIndex: number, toIndex: number) => void;
  children: React.ReactNode;
  draggable?: boolean; 
}

export const DraggableCard: React.FC<DraggableCardProps> = ({ id, index, moveCard, children, draggable = true }) => {
  const { oemColor } = useContext(AppConfigContext);
  const [isHovering, setIsHovering] = useState(false);
  
  const [{ isDragging }, drag] = useDrag({
    type: ItemTypes.CARD,
    item: { id, index },
    collect: (monitor) => ({
      isDragging: monitor.isDragging(),
    }),
    canDrag: draggable,
  });

  const [, drop] = useDrop({
    accept: ItemTypes.CARD,
    hover: (draggedItem: { id: string | number; index: number }) => {
      // 只有当draggable为true时才允许放置
      if (draggable && draggedItem.index !== index) {
        // 调用父组件传递的移动函数
        moveCard(draggedItem.index, index);
        // 更新拖拽项的索引，确保在连续移动时正确计算位置
        draggedItem.index = index;
      }
    },
  });

  // 创建两个ref：一个用于拖拽图标(应用drag)，一个用于整个卡片(应用drop)
  const dragIconRef = React.useRef<HTMLDivElement>(null);
  const cardRef = React.useRef<HTMLDivElement>(null);
  
  // 应用drag到拖拽图标，drop到整个卡片
  drag(dragIconRef);
  drop(cardRef);

  return (
    <div
      ref={cardRef}
      onMouseEnter={() => setIsHovering(true)}
      onMouseLeave={() => setIsHovering(false)}
      style={{
        opacity: isDragging ? 0.5 : 1,
        backgroundColor: isHovering ? oemColor.colorPrimaryBg : '#fff',
        transition: 'opacity 0.2s ease, background-color 0.2s ease, transform 0.1s ease',
        display: 'flex',
        alignItems: 'center',
        transform: isDragging ? 'scale(1.02)' : 'scale(1)',
      }}
    >
        <div 
          ref={dragIconRef}
          style={{
            cursor: draggable ? 'grab' : 'not-allowed',
            transition: 'opacity 0.2s ease',
            minWidth: '20px',
            textAlign: 'center',
            color: '#666',
            display: 'flex',
            alignItems: 'center',
            padding: '0 8px',
          }}
          onMouseDown={(e) => {
            if(draggable) {
              e.stopPropagation();
              e.currentTarget.style.cursor = 'grabbing';
            }
          }}
          onMouseUp={(e) => {
            if(draggable) {
              e.stopPropagation();
              e.currentTarget.style.cursor = 'grab';
            }
          }}
        >
          <DragIcon style={{ width: '16px', height: '16px' }} />
        </div>
      <div style={{ flex: 1 }}>{children}</div>
    </div>
  );
};
