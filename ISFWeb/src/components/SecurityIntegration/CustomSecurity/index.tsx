import React, { useState, useEffect, useMemo } from 'react';
import intl from 'react-intl-universal';
import { DndProvider } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';
import { Modal, Button, message, Input, Empty } from 'antd';
import { trim } from 'lodash';
import { DraggableCard } from '../../DraggableCard';
import { customedSecuName } from '../helper';
import AddIcon from "../../../icons/add.svg";
import EditIcon from "../../../icons/edit.svg";
import DeleteIcon from "../../../icons/delete.svg";
import SaveIcon from "../../../icons/save.svg";
import CloseIcon from "../../../icons/close.svg";
import EmptyIcon from "../../../icons/empty.png";
import styles from './styles';
import __ from './locale';

interface CustomSecurityProps {
    customSecurity: any[];
    hasFileCsfInit: boolean;
    hasUserCsfInit: boolean;
    triggerConfirmCust: (customSecurInfo: any[]) => void;
    closeCustomSecurity: () => void;
}

interface UserLevelItem {
    name: string;
    isEditing?: boolean;
    editValue?: string;
}

interface UserLevel {
    type: string;
    label: string;
    value: UserLevelItem[];
}
const CustomSecurity: React.FC<CustomSecurityProps> = ({ tempCsfLevel, triggerConfirmCust, closeCustomSecurity }) => {
    const [userLevels, setUserLevels] = useState<UserLevel[]>([])
    const [editingId, setEditingId] = useState<string | null>(null)
    
    // 检测是否有任何项处于编辑或添加状态
    const hasActiveEditOrAdd = useMemo(() => {
        return userLevels.some(level => 
            level.value.some(item => item.isEditing || item.isAdding)
        );
    }, [userLevels]);
    
    // 检测是否有密级类型的value为空
    const hasEmptyLevelValue = useMemo(() => {
        return userLevels.some(level => level.value.length === 0);
    }, [userLevels]);

    // 处理编辑状态切换
    const handleEdit = (type: string, index: number, itemId: string) => {
        // 设置当前编辑项ID，防止同时编辑多个项目
        setEditingId(itemId);
        
        // 更新对应项的isEditing状态为true，同时将其他所有项的isEditing设置为false
        const updatedLevels = [...userLevels];
        
        // 先将所有项的isEditing设置为false
        updatedLevels.forEach(level => {
            level.value = level.value.map(item => ({
                ...item,
                isEditing: false
            }));
        });
        
        // 然后设置当前项为编辑状态
        const currentLevelType = updatedLevels.find(level => level.type === type);
        
        if (currentLevelType) {
            const updatedValues = [...currentLevelType.value];
            updatedValues[index] = {
                ...updatedValues[index],
                isEditing: true
            };
            currentLevelType.value = updatedValues;
            setUserLevels(updatedLevels);
        }
    };

    // 保存编辑
    const handleSaveEdit = (type: string, index: number, newName: string) => {
        const trimmedName = trim(newName);
        const updatedLevels = [...userLevels];
        const currentLevelType = updatedLevels.find(level => level.type === type);
        
        if (currentLevelType) {
            // 检查名称是否为空
            if (!trimmedName) {
                currentLevelType.errorMessage = __('密级名称不能为空。');
                setUserLevels(updatedLevels);
                return;
            }

            // 检查名称是否重复
            const nameExists = currentLevelType.value.some((val, idx) => 
                idx !== index && val.name === trimmedName
            );
            
            if (nameExists) {
                currentLevelType.errorMessage = __('该密级名称已存在。')
                setUserLevels(updatedLevels);
                return;
            }
            
            // 检查名称是否包含特殊字符
            if (customedSecuName(trimmedName)) {
                currentLevelType.errorMessage = __('密级名称不能包含 / : * ? " < > | 特殊字符。')
                setUserLevels(updatedLevels);
                return;
            }
            
            // 检查名称长度
            if (trimmedName.length > 80) {
                currentLevelType.errorMessage = __('密级名称不能超过80个字符。')
                setUserLevels(updatedLevels);
                return;
            }
            
            // 更新名称并重置所有编辑状态
            setEditingId(null);
            const updatedLevelType = updatedLevels.find(level => level.type === type);
            
            if (updatedLevelType) {
                // 先将所有项的isEditing设置为false
                updatedLevels.forEach(level => {
                    level.value = level.value.map(item => ({
                        ...item,
                        isEditing: false,
                        isAdding: false // 清除添加标记
                    }));
                });
                
                // 然后更新当前项的名称
                const updatedValues = [...updatedLevelType.value];
                updatedValues[index] = {
                    ...updatedValues[index],
                    name: trimmedName,
                    editValue: undefined,
                    isAdding: false // 对于新添加的项，移除isAdding标记
                };
                updatedLevelType.value = updatedValues;
                setUserLevels(updatedLevels);
            }
        }
    };
    
    // 取消编辑
    const handleCancelEdit = (type: string, index: number) => {
        // 重置编辑状态
        setEditingId(null);
        
        const updatedLevels = [...userLevels];
        const updatedLevelType = updatedLevels.find(level => level.type === type);
        
        if (updatedLevelType) {
            const updatedValues = [...updatedLevelType.value];
            
            // 检查当前项是否是正在添加的项
            const isAddingItem = updatedValues[index]?.isAdding === true;
            
            if (isAddingItem) {
                // 如果是正在添加的项，直接从列表中移除
                updatedValues.splice(index, 1);
                updatedLevelType.value = updatedValues;
            } else {
                // 如果是普通编辑项，直接更新该项的状态
                updatedValues[index] = {
                    ...updatedValues[index],
                    isEditing: false,
                    editValue: undefined
                };
                updatedLevelType.value = updatedValues;
            }
            
            // 重置错误消息
            updatedLevelType.errorMessage = '';
            
            setUserLevels(updatedLevels);
        }
    };
    
    // 删除密级项
    const handleDelete = (type: string, index: number) => {
        const updatedLevels = [...userLevels];
        const currentLevelType = updatedLevels.find(level => level.type === type);
        
        if (currentLevelType) {
            
            // 移除指定索引的项
            const updatedValues = currentLevelType.value.filter((_, idx) => idx !== index);
            currentLevelType.value = updatedValues;
            
            // 如果删除的是当前编辑的项，重置编辑状态
            if (editingId === `${type}-${index}`) {
                setEditingId(null);
            }
            
            setUserLevels(updatedLevels);
        }
    };
    
    // 添加密级项
    const handleAddLevel = (type: string) => {
        const updatedLevels = [...userLevels];
        const currentLevelType = updatedLevels.find(level => level.type === type);
        
        if (currentLevelType) {
            // 检查密级数量是否超过限制
            const MAX_LEVELS = 11;
            if (currentLevelType.value.length >= MAX_LEVELS) {
                message.info('最多添加11个用户密级')
                return;
            }
            
            // 检查是否已经有一个处于编辑状态的添加项
            const hasAddingItem = currentLevelType.value.some(item => item.isAdding || item.isEditing);
            if (hasAddingItem) {
                return; // 已经有一个添加项在编辑中，不重复添加
            }
            
            // 创建一个临时的添加项，处于编辑状态
            const tempAddItem = {
                name: '', // 空名称，等待用户输入
                isEditing: true, // 初始为编辑状态
                isAdding: true, // 标记为正在添加的项
                editValue: '' // 初始编辑值为空
            };
            
            // 将临时项添加到列表开头
            const updatedValues = [tempAddItem, ...currentLevelType.value];
            currentLevelType.value = updatedValues;
            
            // 设置当前编辑ID，防止同时编辑多个项目
            setEditingId(`${type}-${tempAddItem.name}-adding`);
            setUserLevels(updatedLevels);
        }
    };

    const handleSave = () => {
        const formatLevels = userLevels.map((cur) => {
            const values = cur.value.map((item) => ({ name: item.name }))
            return {
                ...cur,
                value: values.reverse()
            }
        })
        triggerConfirmCust(formatLevels)
    }

    useEffect(() => {
        const formatUserLevels = tempCsfLevel.map((cur) => {
            const value = cur.value.map((item) => ({name: item.name}))
            return {
                ...cur,
                value: value.reverse()
            }
        })
        setUserLevels(formatUserLevels)
    }, [])

    return  (
        <DndProvider backend={HTML5Backend}>
            <Modal
                centered
                maskClosable={false}
                open={true}
                width={680}
                title={__('自定义用户密级')}
                onCancel={closeCustomSecurity}
                footer={[
                    <Button key="confirm" type="primary" onClick={handleSave} disabled={hasActiveEditOrAdd || hasEmptyLevelValue}>
                        {intl.get('ok')}
                    </Button>,
                    <Button key="cancel" onClick={closeCustomSecurity}>
                        {intl.get('cancel')}
                    </Button>,
                ]}
                getContainer={document.getElementById("isf-web-plugins") as HTMLElement}
            >
                <div className={styles['custom-security']}>
                    <div className={styles['description']}>{__('下方密级列表从上到下，密级【由高到低】排列；拖动列表项可调整密级的高低顺序')}</div>
                    <div className={styles['security-container']}>
                        {
                            userLevels.map((cur, index) => {
                                return (
                                    <div 
                                        key={cur.type} 
                                        style={{ 
                                            width: userLevels.length === 1 ? '100%' : '50%',
                                            boxSizing: 'border-box',
                                            marginRight: index < userLevels.length - 1 ? '20px' : '0'
                                        }}
                                        className={styles['security']}
                                    >
                                        <div className={styles['header']}>
                                            <div className={styles['header-label']}>{cur.label}</div>
                                            <div>
                                                <Button 
                                                    color="primary" 
                                                    variant="link" 
                                                    icon={<AddIcon style={{width: "14px", height: "14px"}}/>} 
                                                    onClick={() => handleAddLevel(cur.type)}
                                                    disabled={hasActiveEditOrAdd}
                                                >
                                                    {__('添加密级')}
                                                </Button>
                                            </div>
                                        </div>
                                        <div className={styles['main']}>
                                            {
                                                cur.value.length ? cur.value.map((item, index) => {
                                                    return (
                                                        <DraggableCard 
                                                            id={`${cur.type}-${item.name}`}
                                                            index={index} 
                                                            draggable={!(item.isEditing || item.isAdding)}
                                                            moveCard={(fromIndex, toIndex) => {
                                                                const updatedLevels = [...userLevels];
                                                                const currentLevelType = updatedLevels.find(level => level.type === cur.type);
                                                                
                                                                if (currentLevelType) {
                                                                    const updatedValues = [...currentLevelType.value];
                                                                    const [movedItem] = updatedValues.splice(fromIndex, 1);
                                                                    updatedValues.splice(toIndex, 0, movedItem);

                                                                    currentLevelType.value = updatedValues;
                                                                    setUserLevels(updatedLevels);
                                                                }
                                                            }}
                                                        >
                                                            <div 
                                                                key={index} 
                                                                className={styles['item']}
                                                            >
                                                                {item.isEditing ? (
                                                                        <div className={styles['item-edit-content']}>
                                                                            <Input 
                                                                                style={{ width: 'calc(100% - 64px)' }} 
                                                                                defaultValue={item.name}
                                                                                placeholder={__('请输入密级名称')}
                                                                                onChange={(e) => {
                                                                                    const updatedLevels = [...userLevels];
                                                                                    const currentLevelType = updatedLevels.find(level => level.type === cur.type);
                                                                                    if (currentLevelType) {
                                                                                        const updatedValues = [...currentLevelType.value];
                                                                                        updatedValues[index] = {
                                                                                            ...updatedValues[index],
                                                                                            editValue: e.target.value
                                                                                        };
                                                                                        currentLevelType.value = updatedValues;
                                                                                        currentLevelType.errorMessage = '';
                                                                                        setUserLevels(updatedLevels);
                                                                                    }
                                                                                }}
                                                                                onPressEnter={(e) => {
                                                                                    const newName = e.target.value || item.name;
                                                                                    handleSaveEdit(cur.type, index, newName);
                                                                                }}
                                                                                autoFocus={item.isAdding}
                                                                            />
                                                                            <Button
                                                                                color="primary"  
                                                                                variant="link"
                                                                                icon={<SaveIcon style={{width: "14px", height: "14px"}} />}  
                                                                                onClick={() => {
                                                                                    const newName = item.isEditing && item.editValue !== undefined ? item.editValue : item.name;
                                                                                    handleSaveEdit(cur.type, index, newName);
                                                                                }}
                                                                            />
                                                                            <Button
                                                                                color="default"  
                                                                                variant="text"
                                                                                icon={<CloseIcon style={{width: "14px", height: "14px"}} />}  
                                                                                onClick={() => handleCancelEdit(cur.type, index)}
                                                                            />
                                                                        </div>
                                                                    ) : (
                                                                            <div
                                                                                className={styles['item-name']}
                                                                                title={item.name}
                                                                            >
                                                                                {item.name}
                                                                            </div>
                                                                    )
                                                                }
                                                                {
                                                                    item.isEditing ? null : (
                                                                        <div>
                                                                            <Button 
                                                                                color="default" 
                                                                                variant="text" 
                                                                                icon={<EditIcon style={{width: "14px", height: "14px"}}/>} 
                                                                                onClick={() => handleEdit(cur.type, index, `${cur.type}-${item.name}`)}
                                                                                disabled={editingId}
                                                                            />
                                                                            <Button 
                                                                                color="default" 
                                                                                variant="text" 
                                                                                icon={<DeleteIcon style={{width: "14px", height: "14px"}}/>} 
                                                                                onClick={() => {
                                                                                    handleDelete(cur.type, index);
                                                                                }}
                                                                            />
                                                                        </div>
                                                                    )
                                                                }
                                                            </div>
                                                        </DraggableCard>
                                                    )
                                                }) : (
                                                    <div className={styles['empty']}>
                                                        <Empty
                                                            image={EmptyIcon}
                                                            description={__('暂无用户密级，请添加')}
                                                        />
                                                    </div>
                                                )
                                            }
                                        </div>
                                        {
                                            cur.errorMessage ?
                                            <div className={styles['error-message']}>
                                                {cur.errorMessage}
                                            </div> : null
                                        }
                                    </div>
                                )
                            })
                        }
                    </div>
                </div>
            </Modal>
        </DndProvider>
    )
}

export default CustomSecurity