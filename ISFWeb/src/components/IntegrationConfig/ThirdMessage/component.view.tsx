import * as React from 'react';
import { Button, Centered } from '@/ui/ui.desktop';
import Card from './Card/component.view';
import ThirdMessageBase from './component.base';
import styles from './styles.view.css';
import __ from './locale';

export default class ThirdMessage extends ThirdMessageBase {
    render() {
        const { thirdPartyConfig, editedCards } = this.state;
        return (
            <div className={styles['container']}>
                {
                    thirdPartyConfig.map((item, index) => (
                        <Card
                            key={item.thirdPartyName}
                            swf={this.props.swf}
                            configInfo={item}
                            showDeleteIcon={index !== 0}
                            onRequestEditedCardIncrease={(indexId) => { this.handleEditedCardIncrease(indexId) }}
                            onRequestEditedCardDecrease={(indexId) => { this.handleEditedCardDecrease(indexId) }}
                            onRequestDeleteUnSavedCard={() => this.handleDeleteUnSavedCard()}
                            onRequestDeleteSavedCard={(indexId) => this.handleDeleteSavedCard(indexId)}
                            onRequestAddConfigSuccess={({ indexId, thirdPartyName, internalConfig, enabled, pluginClassName, messages }) => this.handleAddConfigSuccess({ indexId, thirdPartyName, internalConfig, enabled, pluginClassName, messages })}
                            onRequestEditConfigSuccess={({ indexId, thirdPartyName, internalConfig, enabled, plugin, pluginClassName, messages }) => this.handleEditConfigSuccess({ indexId, thirdPartyName, internalConfig, enabled, plugin, pluginClassName, messages })}
                            onRequestValidateState={this.validate}
                        />
                    ))
                }
                <div className={styles['add-wrap']}>
                    <Centered>
                        <Button
                            icon={'\uf089'}
                            width={'auto'}
                            size={18}
                            className={styles['add-btn']}
                            disabled={(thirdPartyConfig.length && !thirdPartyConfig[thirdPartyConfig.length - 1].indexId) || editedCards.length}
                            onClick={() => this.handleAddCard()}
                        >
                            {__('添加第三方应用')}
                        </Button>
                    </Centered>
                </div>
            </div>
        )
    }
}