package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname UCSIsAegeanNumbers C.xmlUCSIsAegeanNumbers
func UCSIsAegeanNumbers(code c.Int) c.Int

//go:linkname UCSIsAlphabeticPresentationForms C.xmlUCSIsAlphabeticPresentationForms
func UCSIsAlphabeticPresentationForms(code c.Int) c.Int

//go:linkname UCSIsArabic C.xmlUCSIsArabic
func UCSIsArabic(code c.Int) c.Int

//go:linkname UCSIsArabicPresentationFormsA C.xmlUCSIsArabicPresentationFormsA
func UCSIsArabicPresentationFormsA(code c.Int) c.Int

//go:linkname UCSIsArabicPresentationFormsB C.xmlUCSIsArabicPresentationFormsB
func UCSIsArabicPresentationFormsB(code c.Int) c.Int

//go:linkname UCSIsArmenian C.xmlUCSIsArmenian
func UCSIsArmenian(code c.Int) c.Int

//go:linkname UCSIsArrows C.xmlUCSIsArrows
func UCSIsArrows(code c.Int) c.Int

//go:linkname UCSIsBasicLatin C.xmlUCSIsBasicLatin
func UCSIsBasicLatin(code c.Int) c.Int

//go:linkname UCSIsBengali C.xmlUCSIsBengali
func UCSIsBengali(code c.Int) c.Int

//go:linkname UCSIsBlockElements C.xmlUCSIsBlockElements
func UCSIsBlockElements(code c.Int) c.Int

//go:linkname UCSIsBopomofo C.xmlUCSIsBopomofo
func UCSIsBopomofo(code c.Int) c.Int

//go:linkname UCSIsBopomofoExtended C.xmlUCSIsBopomofoExtended
func UCSIsBopomofoExtended(code c.Int) c.Int

//go:linkname UCSIsBoxDrawing C.xmlUCSIsBoxDrawing
func UCSIsBoxDrawing(code c.Int) c.Int

//go:linkname UCSIsBraillePatterns C.xmlUCSIsBraillePatterns
func UCSIsBraillePatterns(code c.Int) c.Int

//go:linkname UCSIsBuhid C.xmlUCSIsBuhid
func UCSIsBuhid(code c.Int) c.Int

//go:linkname UCSIsByzantineMusicalSymbols C.xmlUCSIsByzantineMusicalSymbols
func UCSIsByzantineMusicalSymbols(code c.Int) c.Int

//go:linkname UCSIsCJKCompatibility C.xmlUCSIsCJKCompatibility
func UCSIsCJKCompatibility(code c.Int) c.Int

//go:linkname UCSIsCJKCompatibilityForms C.xmlUCSIsCJKCompatibilityForms
func UCSIsCJKCompatibilityForms(code c.Int) c.Int

//go:linkname UCSIsCJKCompatibilityIdeographs C.xmlUCSIsCJKCompatibilityIdeographs
func UCSIsCJKCompatibilityIdeographs(code c.Int) c.Int

//go:linkname UCSIsCJKCompatibilityIdeographsSupplement C.xmlUCSIsCJKCompatibilityIdeographsSupplement
func UCSIsCJKCompatibilityIdeographsSupplement(code c.Int) c.Int

//go:linkname UCSIsCJKRadicalsSupplement C.xmlUCSIsCJKRadicalsSupplement
func UCSIsCJKRadicalsSupplement(code c.Int) c.Int

//go:linkname UCSIsCJKSymbolsandPunctuation C.xmlUCSIsCJKSymbolsandPunctuation
func UCSIsCJKSymbolsandPunctuation(code c.Int) c.Int

//go:linkname UCSIsCJKUnifiedIdeographs C.xmlUCSIsCJKUnifiedIdeographs
func UCSIsCJKUnifiedIdeographs(code c.Int) c.Int

//go:linkname UCSIsCJKUnifiedIdeographsExtensionA C.xmlUCSIsCJKUnifiedIdeographsExtensionA
func UCSIsCJKUnifiedIdeographsExtensionA(code c.Int) c.Int

//go:linkname UCSIsCJKUnifiedIdeographsExtensionB C.xmlUCSIsCJKUnifiedIdeographsExtensionB
func UCSIsCJKUnifiedIdeographsExtensionB(code c.Int) c.Int

//go:linkname UCSIsCherokee C.xmlUCSIsCherokee
func UCSIsCherokee(code c.Int) c.Int

//go:linkname UCSIsCombiningDiacriticalMarks C.xmlUCSIsCombiningDiacriticalMarks
func UCSIsCombiningDiacriticalMarks(code c.Int) c.Int

//go:linkname UCSIsCombiningDiacriticalMarksforSymbols C.xmlUCSIsCombiningDiacriticalMarksforSymbols
func UCSIsCombiningDiacriticalMarksforSymbols(code c.Int) c.Int

//go:linkname UCSIsCombiningHalfMarks C.xmlUCSIsCombiningHalfMarks
func UCSIsCombiningHalfMarks(code c.Int) c.Int

//go:linkname UCSIsCombiningMarksforSymbols C.xmlUCSIsCombiningMarksforSymbols
func UCSIsCombiningMarksforSymbols(code c.Int) c.Int

//go:linkname UCSIsControlPictures C.xmlUCSIsControlPictures
func UCSIsControlPictures(code c.Int) c.Int

//go:linkname UCSIsCurrencySymbols C.xmlUCSIsCurrencySymbols
func UCSIsCurrencySymbols(code c.Int) c.Int

//go:linkname UCSIsCypriotSyllabary C.xmlUCSIsCypriotSyllabary
func UCSIsCypriotSyllabary(code c.Int) c.Int

//go:linkname UCSIsCyrillic C.xmlUCSIsCyrillic
func UCSIsCyrillic(code c.Int) c.Int

//go:linkname UCSIsCyrillicSupplement C.xmlUCSIsCyrillicSupplement
func UCSIsCyrillicSupplement(code c.Int) c.Int

//go:linkname UCSIsDeseret C.xmlUCSIsDeseret
func UCSIsDeseret(code c.Int) c.Int

//go:linkname UCSIsDevanagari C.xmlUCSIsDevanagari
func UCSIsDevanagari(code c.Int) c.Int

//go:linkname UCSIsDingbats C.xmlUCSIsDingbats
func UCSIsDingbats(code c.Int) c.Int

//go:linkname UCSIsEnclosedAlphanumerics C.xmlUCSIsEnclosedAlphanumerics
func UCSIsEnclosedAlphanumerics(code c.Int) c.Int

//go:linkname UCSIsEnclosedCJKLettersandMonths C.xmlUCSIsEnclosedCJKLettersandMonths
func UCSIsEnclosedCJKLettersandMonths(code c.Int) c.Int

//go:linkname UCSIsEthiopic C.xmlUCSIsEthiopic
func UCSIsEthiopic(code c.Int) c.Int

//go:linkname UCSIsGeneralPunctuation C.xmlUCSIsGeneralPunctuation
func UCSIsGeneralPunctuation(code c.Int) c.Int

//go:linkname UCSIsGeometricShapes C.xmlUCSIsGeometricShapes
func UCSIsGeometricShapes(code c.Int) c.Int

//go:linkname UCSIsGeorgian C.xmlUCSIsGeorgian
func UCSIsGeorgian(code c.Int) c.Int

//go:linkname UCSIsGothic C.xmlUCSIsGothic
func UCSIsGothic(code c.Int) c.Int

//go:linkname UCSIsGreek C.xmlUCSIsGreek
func UCSIsGreek(code c.Int) c.Int

//go:linkname UCSIsGreekExtended C.xmlUCSIsGreekExtended
func UCSIsGreekExtended(code c.Int) c.Int

//go:linkname UCSIsGreekandCoptic C.xmlUCSIsGreekandCoptic
func UCSIsGreekandCoptic(code c.Int) c.Int

//go:linkname UCSIsGujarati C.xmlUCSIsGujarati
func UCSIsGujarati(code c.Int) c.Int

//go:linkname UCSIsGurmukhi C.xmlUCSIsGurmukhi
func UCSIsGurmukhi(code c.Int) c.Int

//go:linkname UCSIsHalfwidthandFullwidthForms C.xmlUCSIsHalfwidthandFullwidthForms
func UCSIsHalfwidthandFullwidthForms(code c.Int) c.Int

//go:linkname UCSIsHangulCompatibilityJamo C.xmlUCSIsHangulCompatibilityJamo
func UCSIsHangulCompatibilityJamo(code c.Int) c.Int

//go:linkname UCSIsHangulJamo C.xmlUCSIsHangulJamo
func UCSIsHangulJamo(code c.Int) c.Int

//go:linkname UCSIsHangulSyllables C.xmlUCSIsHangulSyllables
func UCSIsHangulSyllables(code c.Int) c.Int

//go:linkname UCSIsHanunoo C.xmlUCSIsHanunoo
func UCSIsHanunoo(code c.Int) c.Int

//go:linkname UCSIsHebrew C.xmlUCSIsHebrew
func UCSIsHebrew(code c.Int) c.Int

//go:linkname UCSIsHighPrivateUseSurrogates C.xmlUCSIsHighPrivateUseSurrogates
func UCSIsHighPrivateUseSurrogates(code c.Int) c.Int

//go:linkname UCSIsHighSurrogates C.xmlUCSIsHighSurrogates
func UCSIsHighSurrogates(code c.Int) c.Int

//go:linkname UCSIsHiragana C.xmlUCSIsHiragana
func UCSIsHiragana(code c.Int) c.Int

//go:linkname UCSIsIPAExtensions C.xmlUCSIsIPAExtensions
func UCSIsIPAExtensions(code c.Int) c.Int

//go:linkname UCSIsIdeographicDescriptionCharacters C.xmlUCSIsIdeographicDescriptionCharacters
func UCSIsIdeographicDescriptionCharacters(code c.Int) c.Int

//go:linkname UCSIsKanbun C.xmlUCSIsKanbun
func UCSIsKanbun(code c.Int) c.Int

//go:linkname UCSIsKangxiRadicals C.xmlUCSIsKangxiRadicals
func UCSIsKangxiRadicals(code c.Int) c.Int

//go:linkname UCSIsKannada C.xmlUCSIsKannada
func UCSIsKannada(code c.Int) c.Int

//go:linkname UCSIsKatakana C.xmlUCSIsKatakana
func UCSIsKatakana(code c.Int) c.Int

//go:linkname UCSIsKatakanaPhoneticExtensions C.xmlUCSIsKatakanaPhoneticExtensions
func UCSIsKatakanaPhoneticExtensions(code c.Int) c.Int

//go:linkname UCSIsKhmer C.xmlUCSIsKhmer
func UCSIsKhmer(code c.Int) c.Int

//go:linkname UCSIsKhmerSymbols C.xmlUCSIsKhmerSymbols
func UCSIsKhmerSymbols(code c.Int) c.Int

//go:linkname UCSIsLao C.xmlUCSIsLao
func UCSIsLao(code c.Int) c.Int

//go:linkname UCSIsLatin1Supplement C.xmlUCSIsLatin1Supplement
func UCSIsLatin1Supplement(code c.Int) c.Int

//go:linkname UCSIsLatinExtendedA C.xmlUCSIsLatinExtendedA
func UCSIsLatinExtendedA(code c.Int) c.Int

//go:linkname UCSIsLatinExtendedB C.xmlUCSIsLatinExtendedB
func UCSIsLatinExtendedB(code c.Int) c.Int

//go:linkname UCSIsLatinExtendedAdditional C.xmlUCSIsLatinExtendedAdditional
func UCSIsLatinExtendedAdditional(code c.Int) c.Int

//go:linkname UCSIsLetterlikeSymbols C.xmlUCSIsLetterlikeSymbols
func UCSIsLetterlikeSymbols(code c.Int) c.Int

//go:linkname UCSIsLimbu C.xmlUCSIsLimbu
func UCSIsLimbu(code c.Int) c.Int

//go:linkname UCSIsLinearBIdeograms C.xmlUCSIsLinearBIdeograms
func UCSIsLinearBIdeograms(code c.Int) c.Int

//go:linkname UCSIsLinearBSyllabary C.xmlUCSIsLinearBSyllabary
func UCSIsLinearBSyllabary(code c.Int) c.Int

//go:linkname UCSIsLowSurrogates C.xmlUCSIsLowSurrogates
func UCSIsLowSurrogates(code c.Int) c.Int

//go:linkname UCSIsMalayalam C.xmlUCSIsMalayalam
func UCSIsMalayalam(code c.Int) c.Int

//go:linkname UCSIsMathematicalAlphanumericSymbols C.xmlUCSIsMathematicalAlphanumericSymbols
func UCSIsMathematicalAlphanumericSymbols(code c.Int) c.Int

//go:linkname UCSIsMathematicalOperators C.xmlUCSIsMathematicalOperators
func UCSIsMathematicalOperators(code c.Int) c.Int

//go:linkname UCSIsMiscellaneousMathematicalSymbolsA C.xmlUCSIsMiscellaneousMathematicalSymbolsA
func UCSIsMiscellaneousMathematicalSymbolsA(code c.Int) c.Int

//go:linkname UCSIsMiscellaneousMathematicalSymbolsB C.xmlUCSIsMiscellaneousMathematicalSymbolsB
func UCSIsMiscellaneousMathematicalSymbolsB(code c.Int) c.Int

//go:linkname UCSIsMiscellaneousSymbols C.xmlUCSIsMiscellaneousSymbols
func UCSIsMiscellaneousSymbols(code c.Int) c.Int

//go:linkname UCSIsMiscellaneousSymbolsandArrows C.xmlUCSIsMiscellaneousSymbolsandArrows
func UCSIsMiscellaneousSymbolsandArrows(code c.Int) c.Int

//go:linkname UCSIsMiscellaneousTechnical C.xmlUCSIsMiscellaneousTechnical
func UCSIsMiscellaneousTechnical(code c.Int) c.Int

//go:linkname UCSIsMongolian C.xmlUCSIsMongolian
func UCSIsMongolian(code c.Int) c.Int

//go:linkname UCSIsMusicalSymbols C.xmlUCSIsMusicalSymbols
func UCSIsMusicalSymbols(code c.Int) c.Int

//go:linkname UCSIsMyanmar C.xmlUCSIsMyanmar
func UCSIsMyanmar(code c.Int) c.Int

//go:linkname UCSIsNumberForms C.xmlUCSIsNumberForms
func UCSIsNumberForms(code c.Int) c.Int

//go:linkname UCSIsOgham C.xmlUCSIsOgham
func UCSIsOgham(code c.Int) c.Int

//go:linkname UCSIsOldItalic C.xmlUCSIsOldItalic
func UCSIsOldItalic(code c.Int) c.Int

//go:linkname UCSIsOpticalCharacterRecognition C.xmlUCSIsOpticalCharacterRecognition
func UCSIsOpticalCharacterRecognition(code c.Int) c.Int

//go:linkname UCSIsOriya C.xmlUCSIsOriya
func UCSIsOriya(code c.Int) c.Int

//go:linkname UCSIsOsmanya C.xmlUCSIsOsmanya
func UCSIsOsmanya(code c.Int) c.Int

//go:linkname UCSIsPhoneticExtensions C.xmlUCSIsPhoneticExtensions
func UCSIsPhoneticExtensions(code c.Int) c.Int

//go:linkname UCSIsPrivateUse C.xmlUCSIsPrivateUse
func UCSIsPrivateUse(code c.Int) c.Int

//go:linkname UCSIsPrivateUseArea C.xmlUCSIsPrivateUseArea
func UCSIsPrivateUseArea(code c.Int) c.Int

//go:linkname UCSIsRunic C.xmlUCSIsRunic
func UCSIsRunic(code c.Int) c.Int

//go:linkname UCSIsShavian C.xmlUCSIsShavian
func UCSIsShavian(code c.Int) c.Int

//go:linkname UCSIsSinhala C.xmlUCSIsSinhala
func UCSIsSinhala(code c.Int) c.Int

//go:linkname UCSIsSmallFormVariants C.xmlUCSIsSmallFormVariants
func UCSIsSmallFormVariants(code c.Int) c.Int

//go:linkname UCSIsSpacingModifierLetters C.xmlUCSIsSpacingModifierLetters
func UCSIsSpacingModifierLetters(code c.Int) c.Int

//go:linkname UCSIsSpecials C.xmlUCSIsSpecials
func UCSIsSpecials(code c.Int) c.Int

//go:linkname UCSIsSuperscriptsandSubscripts C.xmlUCSIsSuperscriptsandSubscripts
func UCSIsSuperscriptsandSubscripts(code c.Int) c.Int

//go:linkname UCSIsSupplementalArrowsA C.xmlUCSIsSupplementalArrowsA
func UCSIsSupplementalArrowsA(code c.Int) c.Int

//go:linkname UCSIsSupplementalArrowsB C.xmlUCSIsSupplementalArrowsB
func UCSIsSupplementalArrowsB(code c.Int) c.Int

//go:linkname UCSIsSupplementalMathematicalOperators C.xmlUCSIsSupplementalMathematicalOperators
func UCSIsSupplementalMathematicalOperators(code c.Int) c.Int

//go:linkname UCSIsSupplementaryPrivateUseAreaA C.xmlUCSIsSupplementaryPrivateUseAreaA
func UCSIsSupplementaryPrivateUseAreaA(code c.Int) c.Int

//go:linkname UCSIsSupplementaryPrivateUseAreaB C.xmlUCSIsSupplementaryPrivateUseAreaB
func UCSIsSupplementaryPrivateUseAreaB(code c.Int) c.Int

//go:linkname UCSIsSyriac C.xmlUCSIsSyriac
func UCSIsSyriac(code c.Int) c.Int

//go:linkname UCSIsTagalog C.xmlUCSIsTagalog
func UCSIsTagalog(code c.Int) c.Int

//go:linkname UCSIsTagbanwa C.xmlUCSIsTagbanwa
func UCSIsTagbanwa(code c.Int) c.Int

//go:linkname UCSIsTags C.xmlUCSIsTags
func UCSIsTags(code c.Int) c.Int

//go:linkname UCSIsTaiLe C.xmlUCSIsTaiLe
func UCSIsTaiLe(code c.Int) c.Int

//go:linkname UCSIsTaiXuanJingSymbols C.xmlUCSIsTaiXuanJingSymbols
func UCSIsTaiXuanJingSymbols(code c.Int) c.Int

//go:linkname UCSIsTamil C.xmlUCSIsTamil
func UCSIsTamil(code c.Int) c.Int

//go:linkname UCSIsTelugu C.xmlUCSIsTelugu
func UCSIsTelugu(code c.Int) c.Int

//go:linkname UCSIsThaana C.xmlUCSIsThaana
func UCSIsThaana(code c.Int) c.Int

//go:linkname UCSIsThai C.xmlUCSIsThai
func UCSIsThai(code c.Int) c.Int

//go:linkname UCSIsTibetan C.xmlUCSIsTibetan
func UCSIsTibetan(code c.Int) c.Int

//go:linkname UCSIsUgaritic C.xmlUCSIsUgaritic
func UCSIsUgaritic(code c.Int) c.Int

//go:linkname UCSIsUnifiedCanadianAboriginalSyllabics C.xmlUCSIsUnifiedCanadianAboriginalSyllabics
func UCSIsUnifiedCanadianAboriginalSyllabics(code c.Int) c.Int

//go:linkname UCSIsVariationSelectors C.xmlUCSIsVariationSelectors
func UCSIsVariationSelectors(code c.Int) c.Int

//go:linkname UCSIsVariationSelectorsSupplement C.xmlUCSIsVariationSelectorsSupplement
func UCSIsVariationSelectorsSupplement(code c.Int) c.Int

//go:linkname UCSIsYiRadicals C.xmlUCSIsYiRadicals
func UCSIsYiRadicals(code c.Int) c.Int

//go:linkname UCSIsYiSyllables C.xmlUCSIsYiSyllables
func UCSIsYiSyllables(code c.Int) c.Int

//go:linkname UCSIsYijingHexagramSymbols C.xmlUCSIsYijingHexagramSymbols
func UCSIsYijingHexagramSymbols(code c.Int) c.Int

//go:linkname UCSIsBlock C.xmlUCSIsBlock
func UCSIsBlock(code c.Int, block *c.Char) c.Int

//go:linkname UCSIsCatC C.xmlUCSIsCatC
func UCSIsCatC(code c.Int) c.Int

//go:linkname UCSIsCatCc C.xmlUCSIsCatCc
func UCSIsCatCc(code c.Int) c.Int

//go:linkname UCSIsCatCf C.xmlUCSIsCatCf
func UCSIsCatCf(code c.Int) c.Int

//go:linkname UCSIsCatCo C.xmlUCSIsCatCo
func UCSIsCatCo(code c.Int) c.Int

//go:linkname UCSIsCatCs C.xmlUCSIsCatCs
func UCSIsCatCs(code c.Int) c.Int

//go:linkname UCSIsCatL C.xmlUCSIsCatL
func UCSIsCatL(code c.Int) c.Int

//go:linkname UCSIsCatLl C.xmlUCSIsCatLl
func UCSIsCatLl(code c.Int) c.Int

//go:linkname UCSIsCatLm C.xmlUCSIsCatLm
func UCSIsCatLm(code c.Int) c.Int

//go:linkname UCSIsCatLo C.xmlUCSIsCatLo
func UCSIsCatLo(code c.Int) c.Int

//go:linkname UCSIsCatLt C.xmlUCSIsCatLt
func UCSIsCatLt(code c.Int) c.Int

//go:linkname UCSIsCatLu C.xmlUCSIsCatLu
func UCSIsCatLu(code c.Int) c.Int

//go:linkname UCSIsCatM C.xmlUCSIsCatM
func UCSIsCatM(code c.Int) c.Int

//go:linkname UCSIsCatMc C.xmlUCSIsCatMc
func UCSIsCatMc(code c.Int) c.Int

//go:linkname UCSIsCatMe C.xmlUCSIsCatMe
func UCSIsCatMe(code c.Int) c.Int

//go:linkname UCSIsCatMn C.xmlUCSIsCatMn
func UCSIsCatMn(code c.Int) c.Int

//go:linkname UCSIsCatN C.xmlUCSIsCatN
func UCSIsCatN(code c.Int) c.Int

//go:linkname UCSIsCatNd C.xmlUCSIsCatNd
func UCSIsCatNd(code c.Int) c.Int

//go:linkname UCSIsCatNl C.xmlUCSIsCatNl
func UCSIsCatNl(code c.Int) c.Int

//go:linkname UCSIsCatNo C.xmlUCSIsCatNo
func UCSIsCatNo(code c.Int) c.Int

//go:linkname UCSIsCatP C.xmlUCSIsCatP
func UCSIsCatP(code c.Int) c.Int

//go:linkname UCSIsCatPc C.xmlUCSIsCatPc
func UCSIsCatPc(code c.Int) c.Int

//go:linkname UCSIsCatPd C.xmlUCSIsCatPd
func UCSIsCatPd(code c.Int) c.Int

//go:linkname UCSIsCatPe C.xmlUCSIsCatPe
func UCSIsCatPe(code c.Int) c.Int

//go:linkname UCSIsCatPf C.xmlUCSIsCatPf
func UCSIsCatPf(code c.Int) c.Int

//go:linkname UCSIsCatPi C.xmlUCSIsCatPi
func UCSIsCatPi(code c.Int) c.Int

//go:linkname UCSIsCatPo C.xmlUCSIsCatPo
func UCSIsCatPo(code c.Int) c.Int

//go:linkname UCSIsCatPs C.xmlUCSIsCatPs
func UCSIsCatPs(code c.Int) c.Int

//go:linkname UCSIsCatS C.xmlUCSIsCatS
func UCSIsCatS(code c.Int) c.Int

//go:linkname UCSIsCatSc C.xmlUCSIsCatSc
func UCSIsCatSc(code c.Int) c.Int

//go:linkname UCSIsCatSk C.xmlUCSIsCatSk
func UCSIsCatSk(code c.Int) c.Int

//go:linkname UCSIsCatSm C.xmlUCSIsCatSm
func UCSIsCatSm(code c.Int) c.Int

//go:linkname UCSIsCatSo C.xmlUCSIsCatSo
func UCSIsCatSo(code c.Int) c.Int

//go:linkname UCSIsCatZ C.xmlUCSIsCatZ
func UCSIsCatZ(code c.Int) c.Int

//go:linkname UCSIsCatZl C.xmlUCSIsCatZl
func UCSIsCatZl(code c.Int) c.Int

//go:linkname UCSIsCatZp C.xmlUCSIsCatZp
func UCSIsCatZp(code c.Int) c.Int

//go:linkname UCSIsCatZs C.xmlUCSIsCatZs
func UCSIsCatZs(code c.Int) c.Int

//go:linkname UCSIsCat C.xmlUCSIsCat
func UCSIsCat(code c.Int, cat *c.Char) c.Int
