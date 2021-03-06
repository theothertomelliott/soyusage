package soyusage_test

import "testing"

// TestAnalyzeConstantMapAccess executes a set of tests to verify the Analyze function's handling
// of using constant values to access map entries.
//
// This allows for situations where there are a fixed set of map fields being accessed, but there is logic
// in the template that selects between them based on other inputs.
func TestAnalyzeConstantMapAccess(t *testing.T) {
	var tests = []analyzeTest{
		{
			name: "handles mapping of string values",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				*/
				{template .main}
					{let $textField}
						c_lifeAbout
					{/let}
					{let $textField2: 'c_other'/}
					{$profile[$textField]}
					{$profile[$textField2]}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"profile": map[string]interface{}{
					"c_other":     "*",
					"c_lifeAbout": "*",
				},
			},
		},
		{
			name: "handles mapping with print directive",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				*/
				{template .main}
					
					{let $textField}
						{'c_lifeAbout' | noAutoescape}
					{/let}
					{let $textField2: 'c_other'/}
					{$profile[$textField]}
					{$profile[$textField2]}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"profile": map[string]interface{}{
					"c_other":     "*",
					"c_lifeAbout": "*",
				},
			},
		},
		{
			name: "handles combined constant and variable values",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				* @param locale
				* @param alternative
				*/
				{template .main}
					{let $textField}
						{if $locale == 'en'}
							c_lifeAbout
						{else}
							{$alternative}
						{/if}
					{/let}
					{$profile[$textField]}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"alternative": "*",
				"locale":      "*",
				"profile": map[string]interface{}{
					"[?]":         "*",
					"c_lifeAbout": "*",
				},
			},
		},
		{
			name: "handles indirect mapping via print and assignment",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				*/
				{template .main}
					{let $textField}
						c_lifeAbout
					{/let}
					{let $indirect}
						{$textField}
					{/let}
					{let $textField2: 'c_other'/}
					{let $indirect2: $textField2/}
					{$profile[$indirect]}
					{$profile[$indirect2]}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"profile": map[string]interface{}{
					"c_other":     "*",
					"c_lifeAbout": "*",
				},
			},
		},
		{
			name: "handles mapping of string values with msg",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				*/
				{template .main}
					{let $textField}
						{msg desc="appropriate key for this language"}c_lifeAbout{/msg}
					{/let}
					{let $textField2: 'c_other'/}
					{$profile[$textField]}
					{$profile[$textField2]}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"profile": map[string]interface{}{
					"c_other":     "*",
					"c_lifeAbout": "*",
				},
			},
		},
		{
			name: "handles mapping from a switch statement",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				* @param category
				* @param about
				*/
				{template .main}
					{let $textField}
						{switch $category}
							{case 'Auto'}
								c_autoAbout
							{case 'Home'}
								c_homeAbout
							{case $about}
								c_lifeAbout
						{/switch}
					{/let}
					{if $profile[$textField]}
						{$profile[$textField]}
					{/if}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"category": "*",
				"about":    "*",
				"profile": map[string]interface{}{
					"c_autoAbout": "*",
					"c_homeAbout": "*",
					"c_lifeAbout": "*",
				},
			},
		},
		{
			name: "handles mapping from a list literal",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				*/
				{template .main}
					{let $list: [
						'c_education',
						'c_awards'
					]/}
					{foreach $item in $list}
						{$profile[$item]}
					{/foreach}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"profile": map[string]interface{}{
					"c_education": "*",
					"c_awards":    "*",
				},
			},
		},
		{
			name: "handles map literal inside list",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				*/
				{template .main}
					{let $list: [
						['field': 'c_education'],
						['field': 'c_awards']
					]/}
					{foreach $item in $list}
						{$profile[$item.field]}
					{/foreach}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"profile": map[string]interface{}{
					"c_education": "*",
					"c_awards":    "*",
				},
			},
		},
		{
			name: "handles ranged map literal, no start",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				*/
				{template .main}
					{foreach $i in range(2)}
						{$profile['field' + $i]}
					{/foreach}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"profile": map[string]interface{}{
					"field0": "*",
					"field1": "*",
				},
			},
		},
		{
			name: "handles ranged map literal",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				*/
				{template .main}
					{foreach $i in range(1,3)}
						{$profile['field' + $i]}
					{/foreach}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"profile": map[string]interface{}{
					"field1": "*",
					"field2": "*",
				},
			},
		},
		{
			name: "handles ranged map literal with increment",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				*/
				{template .main}
					{foreach $i in range(2,6,2)}
						{$profile['field' + $i]}
					{/foreach}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"profile": map[string]interface{}{
					"field2": "*",
					"field4": "*",
				},
			},
		},
		{
			name: "handles keyed map literal",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				*/
				{template .main}
					{let $m: [
						'first': 'c_education',
						'second': 'c_awards'
					]/}
					{foreach $i in keys($m)}
						{$profile[$i]}
					{/foreach}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"profile": map[string]interface{}{
					"c_education": "*",
					"c_awards":    "*",
				},
			},
		},
		{
			name: "handles mapping from an if statement",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				* @param category
				*/
				{template .main}
					{let $textField}
						{if $category == 'Auto'}
							c_autoAbout
						{else}
							c_lifeAbout
						{/if}
					{/let}
					{if $profile[$textField]}
						{$profile[$textField]}
					{/if}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"category": "*",
				"profile": map[string]interface{}{
					"c_autoAbout": "*",
					"c_lifeAbout": "*",
				},
			},
		},
		{
			name: "handles mapping from nested statements",
			templates: map[string]string{
				"test.soy": `
				{namespace test}
				/**
				* @param profile
				* @param category
				* @param about
				*/
				{template .main}
					{let $textField}
						{switch ($category ?: '')}
							{case 'Auto'}
								c_autoAbout
							{case 'Home'}
								c_homeAbout
							{default}
								{if $about == 'Life'}
									c_lifeAbout
								{else}
									c_about
								{/if}
						{/switch}
					{/let}
					{let $value: $profile[$textField] /}
					{$value}
				{/template}
			`,
			},
			templateName: "test.main",
			expected: map[string]interface{}{
				"category": "*",
				"about":    "*",
				"profile": map[string]interface{}{
					"c_autoAbout": "*",
					"c_homeAbout": "*",
					"c_lifeAbout": "*",
					"c_about":     "*",
				},
			},
		},
	}
	testAnalyze(t, tests)
}
