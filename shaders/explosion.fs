#define NUM_PARTICLES	1000
//this
// https://www.shadertoy.com/view/4lfXRf

//other
//https://www.shadertoy.com/view/lsySzd
vec3 pow3(vec3 v, float p)
{
    return pow(abs(v), vec3(p));
}

vec2 noise(vec2 tc)
{
    return (2.*texture(iChannel0, tc).xy-1.).xy;
}

vec3 fireworks(vec2 p)
{
    vec3 color = vec3(0., 0., 0.);
    
        vec2 pos = noise(vec2(0.82, 0.11)*float(fw))*1.5;
    	float time = mod(iTime*3., 6.*(1.+noise(vec2(0.123, 0.987)*float(fw)).x));
        for(int i = 0; i < NUM_PARTICLES; i++)
    	{
        	vec2 dir = noise(vec2(0.512, 0.133)*float(i));
            
            float term = 1./length(p-pos-dir*time)/450.;
            color += pow3(vec3(
                term * 0.9,
                term*0.1,0),
                          1.9);
        }
    return color;
}



void mainImage( out vec4 fragColor, in vec2 fragCoord )
{
	vec2 p = 2. * fragCoord / iResolution.xy - 1.;
    p.x *= iResolution.x / iResolution.y;
    
    vec3 color = fireworks(p);
    fragColor = vec4(color, 1.);
}