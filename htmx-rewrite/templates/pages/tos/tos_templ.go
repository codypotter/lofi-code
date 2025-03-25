// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.856
package tos

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "loficode/templates/components/page"

func TermsOfService() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = page.Page(
			page.HeadConfig{
				Title:       "TOS | loficode",
				Description: "Terms of Service",
				OgImage:     "/assets/images/logo-white.svg",
			},
			[]page.Breadcrumb{},
			content(),
		).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func content() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<section class=\"section\"><div class=\"container content\"><h1 class=\"title is-1\">Terms of Service</h1><p>Welcome to <em>lofi<strong>code</strong></em>! These terms of service govern your use of our website, including all features and services provided. I, Cody Potter, am the sole owner and operator of <em>lofi<strong>code</strong></em>, and I reserve the right to modify or update these terms at any time without prior notice. Please review these terms carefully before using the website. I am not a lawyer, but I have done my best to make these terms clear and easy to understand. If you have any questions or concerns, please contact me at me&#64;codypotter.com.</p><p>I am providing free content on this website in good faith and hope that you will use it responsibly. By using this website, you agree to be bound by these terms of service. If you do not agree with any part of these terms, you may not access the website.</p><p>Misuse of this website may result in termination of your account and legal action. Please be respectful and use this website responsibly.</p><h2 class=\"title is-2\">Acceptance of Terms</h2><p>By accessing or using our website, you agree to be bound by these terms of service. If you do not agree with any part of these terms, you may not access the website.</p><h2 class=\"title is-2\">No Brainers</h2><ul><li>Don't be a jerk</li><li>Don't use this service for anything illegal</li><li>Don't use this service to send spam</li><li>Don't use this service to harm others</li><li>Don't use this service to violate anyone's privacy</li><li>Don't use this service to impersonate someone else</li><li>Don't use this service to cheat or deceive others</li><li>Don't use this service to spread hate speech or discrimination</li><li>Don't use this service to infringe on anyone's intellectual property rights</li><li>Don't use this service to exploit children</li><li>Don't use this service to promote violence</li><li>Don't use this service to promote self-harm</li></ul><h2 class=\"title is-2\">Use of the Website</h2><p>You agree to use the website only for lawful purposes and in a manner consistent with all applicable laws and regulations. You agree not to engage in any activity that disrupts or interferes with the functioning of the website or its services. Violation of these terms may result in termination of your account and legal action.</p><h2 class=\"title is-2\">Intellectional Property</h2><p>All content on the website, including text, graphics, logos, images, and software, is the property of <em>lofi<strong>code</strong></em> (Cody Potter) or its licensors and is protected by copyright and other intellectual property laws. You may not reproduce, distribute, or modify any content from the website without prior written consent.</p><h2 class=\"title is-2\">Limitation of Liability</h2><p>We strive to provide accurate and up-to-date information on the website, but we make no warranties or representations regarding the accuracy or completeness of any content. In no event shall <em>lofi<strong>code</strong></em> be liable for any direct, indirect, incidental, special, or consequential damages arising out of or in any way connected with your use of the website.</p><h2 class=\"title is-2\">Changes to the Terms</h2><p>We reserve the right to modify or update these terms of service at any time without prior notice. Your continued use of the website after any changes constitutes acceptance of those changes.</p><h2 class=\"title is-2\">Governing Law</h2><p>These terms of service shall be governed by and construed in accordance with the laws of North Carolina. Any disputes arising under these terms shall be resolved exclusively in the courts of North Carolina.</p></div></section>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
