#import <UserNotifications/UserNotifications.h>
#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

void SetDarwinAppIcon(const void *data, int len) {
    @autoreleasepool {
        NSData *imgData = [NSData dataWithBytes:data length:len];
        NSImage *icon = [[NSImage alloc] initWithData:imgData];
        if (icon) {
            [NSApp setApplicationIconImage:icon];
        }
    }
}

void SendDarwinNotification(const char *title, const char *message) {
    @autoreleasepool {
        UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];

        // Request authorization (no-op after first grant).
        [center requestAuthorizationWithOptions:(UNAuthorizationOptionAlert | UNAuthorizationOptionSound)
                             completionHandler:^(BOOL granted, NSError *error) {}];

        UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
        content.title = [NSString stringWithUTF8String:title];
        content.body  = [NSString stringWithUTF8String:message];
        content.sound = [UNNotificationSound defaultSound];

        NSString *identifier = [[NSUUID UUID] UUIDString];
        UNNotificationRequest *request =
            [UNNotificationRequest requestWithIdentifier:identifier
                                                 content:content
                                                 trigger:nil];

        [center addNotificationRequest:request withCompletionHandler:^(NSError *error) {
            if (error) {
                NSLog(@"Watchfire notification error: %@", error);
            }
        }];
    }
}
