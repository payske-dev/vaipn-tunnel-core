/*
 * Copyright (c) 2021, Vaipn Inc.
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

#import <XCTest/XCTest.h>

#import "VaipnTunnel.h"


@interface VaipnTunnelDelegate : NSObject <TunneledAppDelegate>
@end
@implementation VaipnTunnelDelegate

- (NSString * _Nullable)getVaipnConfig {
    return @"";
}

@end


@interface VaipnTunnelTests : XCTestCase
@property VaipnTunnelDelegate *vaipnTunnelDelegate;
@end

@implementation VaipnTunnelTests

- (void)setUp {
    [super setUp];
    // Put setup code here. This method is called before the invocation of each test method in the class.
    self.vaipnTunnelDelegate = [[VaipnTunnelDelegate alloc] init];
}

- (void)tearDown {
    // Put teardown code here. This method is called after the invocation of each test method in the class.
    [super tearDown];
}

- (void)testExample {
    // This is an example of a functional test case.
    // Use XCTAssert and related functions to verify your tests produce the correct results.
    
    VaipnTunnel *tunnel = [VaipnTunnel newVaipnTunnel:self.vaipnTunnelDelegate];
    XCTAssertNotNil(tunnel);
}

- (void)testPerformanceExample {
    // This is an example of a performance test case.
    [self measureBlock:^{
        // Put the code you want to measure the time of here.
    }];
}

@end

