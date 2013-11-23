// -*- objc -*-

#include <stdlib.h>
#include <OpenGL/gl.h>

#import "ebiten_controller.h"
#import "ebiten_window.h"

void StartApplication() {
  EbitenController* controller = [[EbitenController alloc] init];
  NSApplication* app = [NSApplication sharedApplication];
  [app setActivationPolicy:NSApplicationActivationPolicyRegular];
  [app setDelegate:controller];
  [app finishLaunching];
  [app activateIgnoringOtherApps:YES];
}

void* CreateGLContext() {
  NSOpenGLPixelFormatAttribute attributes[] = {
    NSOpenGLPFAWindow,
    NSOpenGLPFADoubleBuffer,
    NSOpenGLPFAAccelerated,
    NSOpenGLPFADepthSize, 32,
    0,
  };
  NSOpenGLPixelFormat* format = [[NSOpenGLPixelFormat alloc]
                                  initWithAttributes:attributes];
  NSOpenGLContext* glContext = [[NSOpenGLContext alloc] initWithFormat:format
                                                          shareContext:nil];
  [format release];
  return glContext;
}

void* CreateWindow(size_t width, size_t height, const char* title) {
  NSOpenGLContext* glContext = CreateGLContext();
  [glContext makeCurrentContext];

  NSSize size = NSMakeSize(width, height);
  EbitenWindow* window = [[EbitenWindow alloc]
                            initWithSize:size
                               glContext:glContext];
  [window setTitle: [[NSString alloc] initWithUTF8String:title]];
  [window makeKeyAndOrderFront:nil];

  [glContext setView:[window contentView]];

  return window;
}

void PollEvents(void) {
  for (;;) {
    NSEvent* event = [NSApp nextEventMatchingMask:NSAnyEventMask
                                        untilDate:[NSDate distantPast]
                                           inMode:NSDefaultRunLoopMode
                                          dequeue:YES];
    if (event == nil) {
      break;
    }
    [NSApp sendEvent:event];
  }
}

void BeginDrawing(void* window) {
  glClear(GL_COLOR_BUFFER_BIT);
}

void EndDrawing(void* window) {
  [[(EbitenWindow*)window glContext] flushBuffer];
}