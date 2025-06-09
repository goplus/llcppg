package cargs

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

/**
 * An option is used to describe a flag/argument option submitted when the
 * program is run.
 */

type Option struct {
	Identifier    c.Char
	AccessLetters *c.Char
	AccessName    *c.Char
	ValueName     *c.Char
	Description   *c.Char
}

/**
 * A context is used to iterate over all options provided. It stores the parsing
 * state.
 */

type OptionContext struct {
	Options     *Option
	OptionCount c.SizeT
	Argc        c.Int
	Argv        **c.Char
	Index       c.Int
	InnerIndex  c.Int
	ErrorIndex  c.Int
	ErrorLetter c.Char
	ForcedEnd   bool
	Identifier  c.Char
	Value       *c.Char
}

// llgo:type C
type Printer func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{}) c.Int

/**
 * @brief Prepare argument options context for parsing.
 *
 * This function prepares the context for iteration and initializes the context
 * with the supplied options and arguments. After the context has been prepared,
 * it can be used to fetch arguments from it.
 *
 * @param context The context which will be initialized.
 * @param options The registered options which are available for the program.
 * @param option_count The amount of options which are available for the
 * program.
 * @param argc The amount of arguments the user supplied in the main function.
 * @param argv A pointer to the arguments of the main function.
 */
// llgo:link (*OptionContext).OptionInit C.cag_option_init
func (recv_ *OptionContext) OptionInit(options *Option, option_count c.SizeT, argc c.Int, argv **c.Char) {
}

/**
 * @brief Fetches an option from the argument list.
 *
 * This function fetches a single option from the argument list. The context
 * will be moved to that item. Information can be extracted from the context
 * after the item has been fetched.
 * The arguments will be re-ordered, which means that non-option arguments will
 * be moved to the end of the argument list. After all options have been
 * fetched, all non-option arguments will be positioned after the index of
 * the context.
 *
 * @param context The context from which we will fetch the option.
 * @return Returns true if there was another option or false if the end is
 * reached.
 */
// llgo:link (*OptionContext).OptionFetch C.cag_option_fetch
func (recv_ *OptionContext) OptionFetch() bool {
	return false
}

/**
 * @brief Gets the identifier of the option.
 *
 * This function gets the identifier of the option, which should be unique to
 * this option and can be used to determine what kind of option this is.
 *
 * @param context The context from which the option was fetched.
 * @return Returns the identifier of the option.
 */
// llgo:link (*OptionContext).OptionGetIdentifier C.cag_option_get_identifier
func (recv_ *OptionContext) OptionGetIdentifier() c.Char {
	return 0
}

/**
 * @brief Gets the value from the option.
 *
 * This function gets the value from the option, if any. If the option does not
 * contain a value, this function will return NULL.
 *
 * @param context The context from which the option was fetched.
 * @return Returns a pointer to the value or NULL if there is no value.
 */
// llgo:link (*OptionContext).OptionGetValue C.cag_option_get_value
func (recv_ *OptionContext) OptionGetValue() *c.Char {
	return nil
}

/**
 * @brief Gets the current index of the context.
 *
 * This function gets the index within the argv arguments of the context. The
 * context always points to the next item which it will inspect. This is
 * particularly useful to inspect the original argument array, or to get
 * non-option arguments after option fetching has finished.
 *
 * @param context The context from which the option was fetched.
 * @return Returns the current index of the context.
 */
// llgo:link (*OptionContext).OptionGetIndex C.cag_option_get_index
func (recv_ *OptionContext) OptionGetIndex() c.Int {
	return 0
}

/**
 * @brief Retrieves the index of an invalid option.
 *
 * This function retrieves the index of an invalid option if the provided option
 * does not match any of the options specified in the `cag_option` list. This is
 * particularly useful when detailed information about an invalid option is
 * required.
 *
 * @param context Pointer to the context from which the option was fetched.
 * @return Returns the index of the invalid option, or -1 if it is not invalid.
 */
// llgo:link (*OptionContext).OptionGetErrorIndex C.cag_option_get_error_index
func (recv_ *OptionContext) OptionGetErrorIndex() c.Int {
	return 0
}

/**
 * @brief Retrieves the letter character of the invalid option.
 *
 * This function retrieves the character of the invalid option character
 * if the provided option does not match any of the options specified in the
 * `cag_option` list.
 *
 * @param context Pointer to the context from which the option was fetched.
 * @return Returns the letter that was unknown, or 0 otherwise.
 */
// llgo:link (*OptionContext).OptionGetErrorLetter C.cag_option_get_error_letter
func (recv_ *OptionContext) OptionGetErrorLetter() c.Char {
	return 0
}

/**
 * @brief Prints the error associated with the invalid option to the specified
 * destination.
 *
 * This function prints information about the error associated with the invalid
 * option to the specified destination (such as a file stream). It helps in
 * displaying the error of the current context.
 *
 * @param context Pointer to the context from which the option was fetched.
 * @param destination Pointer to the file stream where the error information
 * will be printed.
 */
// llgo:link (*OptionContext).OptionPrintError C.cag_option_print_error
func (recv_ *OptionContext) OptionPrintError(destination *c.FILE) {
}

/**
 * @brief Prints the error associated with the invalid option using user
 * callback.
 *
 * This function prints information about the error associated with the invalid
 * option using user callback. Callback prototype is same with fprintf. It helps
 * in displaying the error of the current context.
 *
 * @param context Pointer to the context from which the option was fetched.
 * @param printer The printer callback function. For example fprintf.
 * @param printer_ctx The parameter for printer callback. For example fprintf
 * could use parameter stderr.
 */
// llgo:link (*OptionContext).OptionPrinterError C.cag_option_printer_error
func (recv_ *OptionContext) OptionPrinterError(printer Printer, printer_ctx c.Pointer) {
}

/**
 * @brief Prints all options to the terminal.
 *
 * This function prints all options to the terminal. This can be used to
 * generate the output for a "--help" option.
 *
 * @param options The options which will be printed.
 * @param option_count The option count which will be printed.
 * @param destination The destination where the output will be printed.
 */
// llgo:link (*Option).OptionPrint C.cag_option_print
func (recv_ *Option) OptionPrint(option_count c.SizeT, destination *c.FILE) {
}

/**
 * @brief Prints all options using user callback.
 *
 * This function prints all options using user callback. This can be used to
 * generate the output for a "--help" option.
 * Using user callback is useful in tiny system without FILE support
 *
 * @param options The options which will be printed.
 * @param option_count The option count which will be printed.
 * @param destination The destination where the output will be printed.
 * @param printer The printer callback function. For example fprintf.
 * @param printer_ctx The parameter for printer callback. For example fprintf
 * could use parameter stderr.
 */
// llgo:link (*Option).OptionPrinter C.cag_option_printer
func (recv_ *Option) OptionPrinter(option_count c.SizeT, printer Printer, printer_ctx c.Pointer) {
}

// llgo:link (*OptionContext).OptionPrepare C.cag_option_prepare
func (recv_ *OptionContext) OptionPrepare(options *Option, option_count c.SizeT, argc c.Int, argv **c.Char) {
}

// llgo:link (*OptionContext).OptionGet C.cag_option_get
func (recv_ *OptionContext) OptionGet() c.Char {
	return 0
}
