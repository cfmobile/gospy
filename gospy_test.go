package gospy_test

import (
	. "github.com/cfmobile/gospy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
	"errors"
)

const (
	kOriginalStringReturn = "original string value"
	kOriginalIntReturn = 12345
	kOriginalFloatReturn = float64(123.45)
	kOriginalBoolReturn = true
)

var kOriginalErrorReturn = errors.New("some error")

var _ = Describe("GoSpy", func() {
	var subject *GoSpy

	var functionToSpy func(string, int, bool) (string, int, float64, bool, error)
	var panicked bool

	BeforeEach(func() {
	    subject = nil
		panicked = false
		functionToSpy = func(string, int, bool) (string, int, float64, bool, error) {
			return kOriginalStringReturn,
				kOriginalIntReturn,
				kOriginalFloatReturn,
				kOriginalBoolReturn,
				kOriginalErrorReturn
		}
	})

	panicRecover := func() {
		panicked = recover() != nil
	}

	Describe("Constructors", func() {

		var constructorSuccessTests = func() {
			It("should not have panicked", func() {
				Expect(panicked).To(BeFalse())
			})

			It("should have returned a valid *GoSpy object", func() {
				Expect(subject).NotTo(BeNil())
			})
		}

		var constructorFailTests = func() {
			It("should have panicked", func() {
				Expect(panicked).To(BeTrue())
			})

			It("should not have returned a valid *GoSpy object", func() {
				Expect(subject).To(BeNil())
			})
		}

		var spyAndFakeSuccessTests = func() {
			It("should have modified the behaviour of the function to return default type values for each of the return values", func() {
				stringResult, intResult, floatResult, boolResult, errorResult := functionToSpy("something", 10, false)

				Expect(stringResult).To(Equal(""))
				Expect(intResult).To(Equal(0))
				Expect(floatResult).To(Equal(0.0))
				Expect(boolResult).To(Equal(false))
				Expect(errorResult).To(BeNil())
			})
		}

	    Describe("Spy", func() {

	        Context("when calling Spy() with a valid function pointer", func() {
				BeforeEach(func() {
					defer panicRecover()
				    subject = Spy(&functionToSpy)
				})

				constructorSuccessTests()

				It("should not have affected the function's behaviour", func() {
					stringResult, intResult, floatResult, boolResult, errorResult := functionToSpy("something", 10, false)

					Expect(stringResult).To(Equal(kOriginalStringReturn))
					Expect(intResult).To(Equal(kOriginalIntReturn))
					Expect(floatResult).To(Equal(kOriginalFloatReturn))
					Expect(boolResult).To(Equal(kOriginalBoolReturn))
					Expect(errorResult).To(Equal(kOriginalErrorReturn))
				})
	        })

			Context("when calling Spy() with a function var", func() {
			    BeforeEach(func() {
			        defer panicRecover()
					subject = Spy(functionToSpy)
			    })

				constructorFailTests()
			})

			Context("when calling Spy() with any other unexpected type", func() {
			    BeforeEach(func() {
			        defer panicRecover()
					someVar := "some random var"
					subject = Spy(&someVar)
			    })

				constructorFailTests()
			})
	    })

		Describe("SpyAndFake", func() {

			Context("when calling SpyAndFake() with a valid function pointer", func() {
			    BeforeEach(func() {
			        defer panicRecover()
					subject = SpyAndFake(&functionToSpy)
			    })

				constructorSuccessTests()

				spyAndFakeSuccessTests()
			})

			Context("when calling SpyAndFake() with a function object", func() {
				BeforeEach(func() {
				    defer panicRecover()
					subject = SpyAndFake(functionToSpy)
				})

				constructorFailTests()
			})

			Context("when calling SpyAndFake() with any other unexpected type", func() {
			    BeforeEach(func() {
			        defer panicRecover()
					someVar := "some random var"
					subject = SpyAndFake(&someVar)
			    })

				constructorFailTests()
			})
		})

		Describe("SpyAndFakeWithReturn", func() {

			Context("when calling SpyAndFakeWithReturn() with a valid function pointer and valid mock return values", func() {
				mockStringValue := "mock value"
				mockIntValue := 1
				mockFloatValue := 2.0
				mockBoolValue := false
				mockErrorValue := errors.New("mock error")

				BeforeEach(func() {
				    defer panicRecover()
					subject = SpyAndFakeWithReturn(&functionToSpy, mockStringValue, mockIntValue, mockFloatValue, mockBoolValue, mockErrorValue)
				})

				constructorSuccessTests()

				It("should have altered the function to just return the mock values specified", func() {
					stringResult, intResult, floatResult, boolResult, errorResult := functionToSpy("something", 10, false)

					Expect(stringResult).To(Equal(mockStringValue))
					Expect(intResult).To(Equal(mockIntValue))
					Expect(floatResult).To(Equal(mockFloatValue))
					Expect(boolResult).To(Equal(mockBoolValue))
					Expect(errorResult).To(Equal(mockErrorValue))
				})
			})

			Context("when calling SpyAndFakeWithReturn() with no fake return values while the monitored function expects return values", func() {
			    BeforeEach(func() {
			        defer panicRecover()
					subject = SpyAndFakeWithReturn(&functionToSpy)
			    })

				// It should succeed and mock the function to return default type values for each of the return values

				constructorSuccessTests()

				spyAndFakeSuccessTests()
			})

			Context("when calling SpyAndFakeWithReturn() with an invalid first argument (not a function pointer)", func() {
			    BeforeEach(func() {
			        defer panicRecover()
					subject = SpyAndFakeWithReturn(functionToSpy, "mock", 1, 2.0, false, nil)
			    })

				constructorFailTests()
			})

			Context("when calling SpyAndFakeWithReturn() with an incorrect number of arguments (not matching the number of return values in the monitored function)", func() {
			    BeforeEach(func() {
			        defer panicRecover()
					subject = SpyAndFakeWithReturn(&functionToSpy, "mock", 1)
			    })

				constructorFailTests()
			})

			Context("when calling SpyAndFakeWithReturn() with an incorrect variable type for any of the mock return values", func() {
			    BeforeEach(func() {
			        defer panicRecover()
					subject = SpyAndFakeWithReturn(&functionToSpy, 0, 1, 2.0, false, nil)
			    })

				constructorFailTests()
			})
		})
	})

	Context("when a valid GoSpy object is created", func() {
		var expectedCalledState bool
		var expectedCallCount int
		var expectedCallList CallList

		// Definition of common tests for each scenario
		var goSpyResetTests = func() {
			Context("when Reset() is called", func() {
				BeforeEach(func() {
					subject.Reset()
				})

				It("should zero the call count", func() {
					Expect(subject.CallCount()).To(BeZero())
				})

				It("should return a nil call list", func() {
					Expect(subject.Calls()).To(BeNil())
				})

				It("should have reset the call indicator", func() {
					Expect(subject.Called()).To(BeFalse())
				})
			})
		}

		var goSpyRestoreTests = func(existingCallCount int, existingCallList CallList) {
			Context("when Restore() is called", func() {
				BeforeEach(func() {
					subject.Restore()
				})

				It("should not have affected the existing call count", func() {
					Expect(subject.CallCount()).To(Equal(existingCallCount))
				})

				It("should not have affected the call list", func() {
					Expect(subject.Calls()).To(Equal(existingCallList))
				})

				It("should no longer monitor subsequent calls to the function", func() {
					Expect(subject.CallCount()).To(Equal(existingCallCount))

					functionToSpy("another call", 101, true)

					Expect(subject.CallCount()).To(Equal(existingCallCount))
					Expect(subject.Calls()).NotTo(ContainElement(ArgList{"another call", 101, true}))
				})
			})
		}

		var goSpyCalledTest = func(expectedCalledState bool) {
			wasCalled := "was"
			if !expectedCalledState {
				wasCalled = "was not"
			}

			It(fmt.Sprintf("should indicate that the function %s Called()", wasCalled), func() {
				Expect(subject.Called()).To(Equal(expectedCalledState))
			})
		}

		var goSpyCallCountTest = func(expectedCallCount int) {
			It(fmt.Sprintf("should indicate a CallCount() of %d", expectedCallCount), func() {
				Expect(subject.CallCount()).To(Equal(expectedCallCount))
			})
		}

		var goSpyCallsTest = func(expectedCallList CallList) {
			msg := "an expected and ordered"
			if expectedCallList == nil {
				msg = "a nil"
			}

			It(fmt.Sprintf("should contain %s list of Calls()", msg), func() {
			    Expect(subject.Calls()).To(Equal(expectedCallList))
			})
		}


		BeforeEach(func() {
			subject = Spy(&functionToSpy)
		})

		Context("as soon as it's created", func() {
			expectedCalledState = false
		    expectedCallCount = 0
			expectedCallList = nil

			goSpyCalledTest(expectedCalledState)

			goSpyCallCountTest(expectedCallCount)

			goSpyCallsTest(expectedCallList)

			goSpyResetTests()

			goSpyRestoreTests(expectedCallCount, expectedCallList)

			Context("when ArgsForCall() is called with no calls in the Spy", func() {
				BeforeEach(func() {
					defer panicRecover()
					subject.ArgsForCall(0)
				})

				It("should panic", func() {
					Expect(panicked).To(BeTrue())
				})
			})
		})

		Context("and the monitored function is called once", func() {
			expectedCalledState = true
			expectedCallCount = 1
			expectedArgList := ArgList{"test value", 101, true}
			expectedCallList = CallList{expectedArgList}

			BeforeEach(func() {
			    functionToSpy("test value", 101, true)
			})

			goSpyCalledTest(expectedCalledState)

			goSpyCallCountTest(expectedCallCount)

			goSpyCallsTest(expectedCallList)

			goSpyResetTests()

			goSpyRestoreTests(expectedCallCount, expectedCallList)

			It("ArgsForCall() should return the arguments that were used in the call", func() {
			    Expect(subject.ArgsForCall(0)).To(Equal(expectedArgList))
			})
		})

		Context("and the monitored function is called several times", func() {
			expectedCalledState = true
			expectedCallCount = 3
			expectedCallList = CallList{
				{"call 1", 1, true},
				{"call 2", 2, false},
				{"call 3", 3, true},
			}

			BeforeEach(func() {
			    functionToSpy("call 1", 1, true)
				functionToSpy("call 2", 2, false)
				functionToSpy("call 3", 3, true)
			})

			goSpyCalledTest(expectedCalledState)

			goSpyCallCountTest(expectedCallCount)

			goSpyCallsTest(expectedCallList)

			goSpyResetTests()

			goSpyRestoreTests(expectedCallCount, expectedCallList)

			It("ArgsForCall(n) should return the arguments for the n-th call (0-based index) ", func() {
			    Expect(subject.ArgsForCall(0)).To(Equal(expectedCallList[0]))
			    Expect(subject.ArgsForCall(1)).To(Equal(expectedCallList[1]))
			    Expect(subject.ArgsForCall(2)).To(Equal(expectedCallList[2]))
			})
		})
	})
})
