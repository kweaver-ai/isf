import * as React from 'react';
import { debounce } from 'lodash';
import classNames from 'classnames';
import ListTipComponent from '../../../../ListTipComponent/component.view';
import { ListTipStatus } from '../../../../ListTipComponent/helper'
import { Select2, PopMenu, SingleSelection, ValidateTip, SweetIcon } from '@/sweet-ui';
import { LazyLoader, SearchInput, Title } from '@/ui/ui.desktop';
import { ValidateMessages, ValidateStatus } from '../../../type';
import { SelectOptionsLimit } from './helper';
import { BlurType, Option, SelectProps } from './type';
import  styles from './style.css';
import __ from './locale';

const { useState, useRef, useEffect } = React;

const { Option: SelectOption } = Select2;

const Select: React.FC<SelectProps> = ({
    selectedKeyOfProps = '',
    inputValueOfProps = '',
    optionsOfProps = [],
    showSearch = false,
    validateStatus,
    loader,
    onChange,
    allowClear = false,
}) => {
    const triggerElement = useRef<any>(null);
    const blurType = useRef<BlurType>(BlurType.Other);
    const count = useRef<number>(0);
    const searchInputRef = useRef<any>(null);
    const lazyLoaderRef = useRef<any>(null);
    const isSearchAction = useRef<boolean>(false);
    const timer = useRef<any>(null);
    const oldInfo = useRef<{
        inputValue: string;
        selectedKey: string | number;
        options: ReadonlyArray<Option>;
    }>({
        inputValue: inputValueOfProps,
        selectedKey: selectedKeyOfProps,
        options: optionsOfProps,
    }); // 存储原有选择信息

    const [options, setOptions] = useState<ReadonlyArray<Option>>(optionsOfProps);
    const [inputValue, setInputValue] = useState<string>(inputValueOfProps);
    const [selectedKey, setSelectedKey] = useState<string | number>(selectedKeyOfProps);
    const [placeholder, setPlaceholder] = useState<string>('');
    const [listTipStatus, setListTipStatus] = useState<ListTipStatus>(ListTipStatus.None);
    const [isFocus, setIsFocus] = useState<boolean>(false);
    const [hover, setHover] = useState<boolean>(false);
    const [arrowHover, setArrowHover] = useState<boolean>(false);

    useEffect(() => {
        // 只有搜索时，才需要执行该逻辑
        if (showSearch) {
            if (isSearchAction.current) {
                if (timer.current) {
                    timer.current.cancel();
                    timer.current = null;
                }

                timer.current = debounce(async () => {
                    setOptions([]);

                    await loadList();

                    lazyLoaderRef.current && lazyLoaderRef.current.reset();

                    timer.current = null;
                }, 200)

                timer.current();
            }
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [inputValue])

    useEffect(() => {
        // 如果options为空，则加载数据
        if (!options.length) {
            loadList();
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const isNotNormal = validateStatus !== ValidateStatus.Normal;

    return (
        <PopMenu
            freeze={false}
            anchorOrigin={['left', 'bottom']}
            alignOrigin={['left', 'top']}
            triggerEvent={showSearch ? 'focus' : 'click'}
            element={[triggerElement.current, 'window']}
            onRequestCloseWhenClick={(close) => close()}
            trigger={({ setPopupVisibleOnFocus, setPopupVisibleOnClick }) =>
                <ValidateTip
                    key='select-trigger'
                    placement={'rightTop'}
                    content={ValidateMessages[validateStatus]}
                    visible={isNotNormal && hover}
                    tipStatus={'error'}
                >
                    <div ref={(ref) => {
                        triggerElement.current = ref;
                    }}>
                        {
                            showSearch ? (
                                <div
                                    className={classNames(
                                        styles['select'],
                                        {
                                            [styles['focus']]: isFocus,
                                            [styles['error']]: isNotNormal,
                                        },
                                    )}
                                    onMouseEnter={() => setHover(true)}
                                    onMouseLeave={() => setHover(false)}
                                >
                                    <div className={styles['input-container']}>
                                        <Title
                                            content={isSearchAction.current ? '' : inputValue}
                                            inline={true}
                                        >
                                            <SearchInput
                                                ref={(ref) => { searchInputRef.current = ref; }}
                                                className={classNames(
                                                    styles['input'],
                                                )}
                                                value={inputValue}
                                                placeholder={placeholder || __('请选择')}
                                                onChange={handleValueChange}
                                                onFocus={() => { handleInputFocus(setPopupVisibleOnFocus) }}
                                                onBlur={handleInputBlur}
                                            />
                                        </Title>
                                    </div>
                                    <div
                                        className={styles['arrow-wrapper']}
                                        onMouseEnter={() => allowClear && inputValue && setArrowHover(true)}
                                        onMouseLeave={() => setArrowHover(false)}
                                        onClick={(e) => {
                                            if (allowClear && inputValue && arrowHover) {
                                                handleClear(e);
                                            } else {
                                                handleInputFocus(setPopupVisibleOnFocus);
                                            }
                                        }}
                                    >
                                        <SweetIcon
                                            role={'sweetui-sweeticon'}
                                            className={styles['arrow']}
                                            name={allowClear && inputValue && arrowHover ? 'x' : 'arrowDown'}
                                            size={16}
                                        />
                                    </div>
                                    {
                                        isNotNormal ?
                                            <SweetIcon
                                                name={'caution'}
                                                size={16}
                                                color={'#e60012'}
                                                className={styles['caution']}
                                                onClick={() => { handleInputFocus(setPopupVisibleOnFocus) }}
                                            /> : null
                                    }
                                </div>
                            ) : (
                                <SingleSelection
                                    placeholder={__('请选择')}
                                    label={inputValue}
                                    width={244}
                                    onClick={setPopupVisibleOnClick}
                                    status={isNotNormal ? 'error' : 'normal'}
                                    onMouseEnter={() => setHover(true)}
                                    onMouseLeave={() => setHover(false)}
                                />
                            )
                        }
                    </div>
                </ValidateTip>
            }
        >
            <div className={styles['container']}>
                <div
                    className={classNames(
                        styles['auto-list'],
                        {
                            [styles['list']]: options.length > 3,
                        },
                    )}
                >
                    {
                        options.length ? (
                            <LazyLoader
                                ref={(ref) => lazyLoaderRef.current = ref}
                                role={'ui-lazyloader'}
                                limit={SelectOptionsLimit}
                                onChange={handleLazyLoader}
                            >
                                {
                                    options.map(({ value_code, value_name }) => {
                                        return (
                                            <SelectOption
                                                value={value_code}
                                                key={value_code}
                                                selected={selectedKey === value_code}
                                                onClick={() => selectOption(value_name, value_code)}
                                            >
                                                {value_name}
                                            </SelectOption>
                                        )
                                    })
                                }
                                {

                                    listTipStatus === ListTipStatus.Loading ? (
                                        <div onClick={(e) => { e.stopPropagation() }}>
                                            <div className={styles['loading-tip-container']}>
                                                <ListTipComponent
                                                    imgSize={24}
                                                    listTipStatus={listTipStatus}
                                                />
                                                <span className={styles['loading']}>{__('加载中...')}</span>
                                            </div>
                                        </div>
                                    ) : null
                                }
                            </LazyLoader>
                        ) : (
                            listTipStatus === ListTipStatus.None ? (
                                <div
                                    className={styles['none']}
                                    onClick={(e) => { e.stopPropagation() }}
                                >
                                    {__(inputValue ? '未找到匹配的结果' : '暂无数据')}
                                </div>
                            ) : null
                        )
                    }
                    {
                        !options.length && listTipStatus === ListTipStatus.Loading ? (
                            <div onClick={(e) => { e.stopPropagation() }}>
                                <div className={styles['loading-tip-container']}>
                                    <ListTipComponent
                                        imgSize={24}
                                        listTipStatus={listTipStatus}
                                    />
                                    <span className={styles['loading']}>{__('加载中...')}</span>
                                </div>
                            </div>
                        ) : null
                    }
                    {
                        listTipStatus === ListTipStatus.LoadFailed ? (
                            <div onClick={(e) => { e.stopPropagation() }}>
                                <div className={styles['load-failed-tip-container']}>
                                    <span>{__('加载失败，')}</span>
                                    <span
                                        className={styles['retry']}
                                        onClick={reloadList}
                                    >
                                        {__('重试')}
                                    </span>
                                </div>
                            </div>
                        ) : null
                    }
                </div>
            </div>
        </PopMenu>
    )

    /**
    * 请求列表数据
    */
    async function loadList(options: ReadonlyArray<Option> = []): Promise<ReadonlyArray<Option>> {
        const offset = options.length;

        if (!count.current || count.current > offset) {
            try {
                setListTipStatus(ListTipStatus.Loading);

                const { entries, total_count } = await loader({ limit: SelectOptionsLimit, offset, searchKey: isSearchAction.current ? inputValue : '' });
                count.current = total_count;

                setOptions([...options, ...entries]);
                setListTipStatus(ListTipStatus.None);

                return [...options, ...entries]
            } catch (error) {
                setOptions([]);
                setListTipStatus(ListTipStatus.LoadFailed);
                return []
            }
        }
        return options
    }

    /**
     * 点击重试，重新加载列表数据
     */
    async function reloadList(e: React.MouseEvent<HTMLSpanElement>): Promise<void> {
        e.stopPropagation();

        await loadList();

        lazyLoaderRef.current && lazyLoaderRef.current.reset();
    }

    /**
     * 懒加载函数
     */
    function handleLazyLoader(): void {
        loadList(options);
    }

    /**
    * 聚焦文本框
    */
    function handleInputFocus(setPopupVisibleOnFocus: () => void): void {
        searchInputRef.current && searchInputRef.current.focus();
        setInputValue('');
        setPlaceholder(inputValue);
        setIsFocus(true);

        setPopupVisibleOnFocus();
    }

    /**
    * 文本框失焦
    */
    function handleInputBlur() {
        setTimeout(() => {
            if (blurType.current === BlurType.Other) {
                setInputValue(placeholder);
                setPlaceholder('');

                // 失焦后，进行了搜索，但是没有选择新的下拉项，恢复原有选择信息
                if (isSearchAction.current && selectedKey) {
                    const { inputValue, selectedKey, options } = oldInfo.current;

                    setOptions(options);

                    updateSelectedInfo(inputValue, selectedKey);
                }
            }

            setIsFocus(false);
            blurType.current = BlurType.Other;
        }, 100)
    }

    /**
     * 选择下拉项
     */
    function selectOption(inputValue: string, selectedKey: string | number): void {
        blurType.current = BlurType.Selected;

        // 保存信息
        oldInfo.current = {
            options,
            inputValue,
            selectedKey,
        }

        updateSelectedInfo(inputValue, selectedKey);
    }

    /**
     * 更新已选信息
     */
    async function updateSelectedInfo(inputValue: string, selectedKey: string | number): Promise<void> {
        setInputValue(inputValue);
        setSelectedKey(selectedKey);
        let optionsTem = options

        if (isSearchAction.current) {
            // 将搜索动作置为false
            isSearchAction.current = false;

            optionsTem = await loadList()
        }

        onChange({
            value_code: selectedKey,
            value_name: inputValue,
        }, optionsTem);
    }

    /**
    * 改变搜索关键字
    */
    function handleValueChange(searchKey: string) {
        // 如果仅仅是聚焦文本框，没有进行输入操作，则不进入该逻辑
        if (!isSearchAction.current && searchKey !== inputValue) {
            isSearchAction.current = true;
        }

        setInputValue(searchKey);
    }

    /**
     * 清除选择
     */
    function handleClear(e: React.MouseEvent) {
        e.stopPropagation();
        setInputValue('');
        setSelectedKey('');
        setPlaceholder('');
        onChange({ value_code: '', value_name: '' }, options);
    }
}

export default Select;