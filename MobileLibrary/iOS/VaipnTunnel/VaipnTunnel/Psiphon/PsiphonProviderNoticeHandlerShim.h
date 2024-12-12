/*
 * Copyright (c) 2020, Vaipn Inc.
 * All rights reserved.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

#import <Foundation/Foundation.h>
#import <Psi/Psi.h>

NS_ASSUME_NONNULL_BEGIN

/// VaipnProviderNoticeHandlerShim provides a shim between GoPsiVaipnProviderNoticeHandler and a block which logs notices.
/// @note This indirection is required because gomobile does not support Objective-C blocks.
@interface VaipnProviderNoticeHandlerShim : NSObject <GoPsiVaipnProviderNoticeHandler>

/// Initialize the notice handler with the given logger.
/// @param logger Logger which will receive notices.
- (id)initWithLogger:(void (^__nonnull)(NSString *_Nonnull))logger;

@end

NS_ASSUME_NONNULL_END
