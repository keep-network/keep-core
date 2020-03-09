import React from 'react'
import * as Icons from './Icons'
import { ClockIndicator } from './Loadable'
import { colors } from '../constants/colors'

export const BANNER_TYPE = {
  SUCCESS: { className: 'success', iconComponent: <Icons.OK color={colors.success} /> },
  PENDING: { className: 'pending', iconComponent: <ClockIndicator /> },
  ERROR: { className: 'error', iconComponent: <Icons.Cross color={colors.error} height={10} width={10} /> },
  DISABLED: { className: 'disabled', iconComponent: null },
}

const Banner = ({
  type,
  title,
  onTitleClick,
  titleClassName,
  subtitle,
  withIcon,
  withCloseIcon,
  onCloseIcon,
  children,
}) => {
  return (
    <div className={`banner banner-${type.className}`}>
      {withIcon &&
        <div className="banner-icon flex">
          {type.iconComponent}
        </div>
      }
      <div className='banner-content-wrapper'>
        <div className={`banner-title ${titleClassName}`} onClick={onTitleClick}>
          {title}
        </div>
        {subtitle &&
            <div className="banner-subtitle">
              {subtitle}
            </div>
        }
      </div>
      { withCloseIcon &&
        <div className='banner-close-icon' onClick={onCloseIcon}>
          <Icons.Cross color={colors[type.className]} height={10} width={10} />
        </div>
      }
      {children}
    </div>
  )
}

Banner.defaultProps = {
  onTitleClick: () => {},
  titleClassName: '',
  withIcon: false,
  withCloseIcon: false,
  onCloseIcon: () => {},
  children: null,
}

export default React.memo(Banner)
