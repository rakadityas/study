# print shape
puts "Raka" #print new line afterward
print "\"Academy\"" # no new llne provided
print " is cool \n" # add new line manually

# ruby works with instruction from each line
# instruction goes through each line

# variables
# ruby will have special container called variable
char_name = "Raka"
char_age = "25"

# option 1 to insert character name
puts "Hello my name is #{char_name}"
puts "My age is #{char_age}"

# option 2 to insert character name
puts ("Hello my name is " + char_name)
puts ("My age is " + char_age)

# data types
string_example = "  Raka   Aditya"
int_example = 10
float_example = 3.2
bool_example = true
nil_example = nil # does not have any value


# handy string lib
string_example.upcase()
string_example.downcase()
string_example.strip()
string_example.length()
string_example.include? "Aditya" # will return true or false value, whether we have the string or not
string_example[2] # get charhcter by index, start from 0
string_example[0, 3] # get the first three charachter
string_example.index("Raka") # get index out of given string
"some string".upcase() # no need to store to variable

# number
some_number = 10
puts "some number: " + some_number.to_s # we can print string and int simultaniously


# arrays
array_names = Array["Raka", "Riki", "Roko"]
array_names = Array[1, "Riki", true] # you can mix it
puts array_names[-1] # using -, will get the index backward
array_names.reverse() # will reverse all element
# array_names.include ? "Raka" # check if current array has element
# array_names.sort() # sort the array

# hash
some_hash = {
    "Jakarta" => "JKT",
    "Malang" => "MLG"
}

# method
def sayHi(name, age=24) # predefined value
    puts "Hello " + name + " " + age.to_s
end
sayHi("Raka") # calling method

def cube(num)
    10
    5 # will return 5, because 5 is the last line of the function
end
puts cube(1)

def cube(num)
    return 10
    5 # will return 10, this line of code won't get executed
end
puts cube(1)

def cube(num)
    return 10, 11
end
puts cube(1)[0] # return multiple values

# for loop
friends = ["Raka", "Riki", "Roko"]

for element in friends
    puts element
end

friends.each do |friend|
    puts friend
end

# handle exception
begin
    num = 10/0
rescue ZeroDivisionError
    puts "something is wrong"
rescue TypeError => e
    puts e
end


# classes and object
class Book #use capital letter; by creating a class, we are creating a custom data type or entity; create a template of new attribute
    attr_accessor :title, :author, :pages #attributes
    def initialize(title, author, pages)
        puts "Creating book" 

        @title = title
        @author = author
        @pages = pages
    end
end

book1 = Book.new("Harry Potter", "JK Rowling", 100) # create new class


class Student
    attr_accessor :name, :major, :gpa
    def initialize(name, major, gpa)
        @name = name
        @major = major
        @gpa = gpa
    end

    def has_honors
        if @gpa >= 3.5 
            return true
        end
        return false
    end
end

student1 = Student.new("Raka", "CS", 4.00)
student2 = Student.new("Riki", "DS", 3.40)

puts student2.has_honors

# inheritance
class Chef 
    def make_chicken
        puts "make chicken"
    end
    def make_salad
        puts "make salad"
    end
    def make_special_dish
        puts "make bbq"
    end
end

class ItalianChef < Chef # inheritance on Ruby works by including < sign
    def make_special_dish
        puts "The chef makes eggplant parm"
    end
    def make_pasta
        puts "The chef makes pasta"
    end
end

chef = Chef.new()
chef.make_chicken


it = ItalianChef.new()
it.make_special_dish
it.make_pasta

# module: store different methods into a collections
module Tools
    def sayhi(name) 
        puts "hello #{name}"
    end
    def saybye(name) 
        puts "bye #{name}"
    end
end

include Tools # need to import first before using it

Tools.sayhi("Raka")
Tools.saybye("Raka")



class Chef 
    def make_chicken
        puts "make chicken"
    end
    def create_salad
        puts "make salad"
    end
    def makeSpecialDish
        puts "make bbq"
    end
    def mk_id_mie
        puts "make indomie"
    end
end

class ItalianChef < Chef 
    def makeSpecialDish
        puts "The chef makes eggplant parm"
    end
    def make_pasta
        puts "The chef makes pasta"
    end
    def make_indomie
        puts "make indomie with meatball"
    end
end



class Chef 
    def make_chicken
        puts "make chicken"
    end
    def make_salad
        puts "make salad"
    end
    def make_special_dish
        puts "make bbq"
    end
end

class ItalianChef < Chef 
    def make_special_dish
        puts "The chef makes eggplant parm"
    end
    def make_pasta
        puts "The chef makes pasta"
    end
end

ada 2



class Animal
	attr_reader :type
	
	def initialize(type)
		@type = type
	end

	def is_a_peacock?
		if type == "peacock"
			return true
		else
			return false
		end
	end
end


class Animal
	attr_reader :type
	
	def initialize(type)
		@type = type
	end

	def is_a_peacock?(animal)
		type == "peacock"
	end
end

1




def is_prime?(n)
	# loop until i * i is greater than n
	while i * i < n
		# check to see if n is evenly divisible by integer i
		if n % i == 0 
			# return false if it is
            puts "[debug] not prime"
			return false
		end
		i += 1
	end
	# return true if n % i == 0 is never true
    puts "[debug] prime"
	true
end

def is_prime?(n)
	# Any factor greater than sqrt(n) has a corresponding factor less than 
	# sqrt(n), so once i >=sqrt(n) we have covered all cases
	while i * i < n
		if n % i == 0 
			return false
		end
		i += 1
	end
	true
end

1



A = :fluorescent
B = :incandescent
C = :led

class Thing
	attr_reader :x
	
	def initialize(x)
		@x = x
	end

	def lightbulb
		if x == A || x == B || x == C
			puts "I am a lightbulb!!"
		else
			puts "I don't know what I am..."
		end
	end
end

class GarageItem
    attr_reader :type

    def initialize(type)
        @type = type
    end

    def am_i_a_lightbulb
        if type == :incandescent || type == :fluorescent || type == :led
            puts "I am a lightbulb!!"
        else
            puts "I don't know what I am."
        end
    end
end


2



class Person

	attr_reader: :height, :hair_color, :dominant_hand, :iq

	def initialize(height, hair_color, dominant_hand, iq)
		@height = height
		@hair_color = hair_color
		@dominant_hand = dominant_hand
		@iq = iq
	end

end


Struct.new("Person", :height, :hair_color, :dominant_hand, :iq)
# or
Person = Struct.new(:height, :hair_color, :dominant_hand, :iq)


sally = Person.new(165, "red", :left, 180, true)


1





puts "What is your major?"
major = gets.chomp

case major
when "Biology"
	puts "Mmm the study of life itself!"
when "Computer Science"
	puts "I'm a computer!"
when "English"
	puts "No way! What's your favorite book?"
when "Math"
	puts "Sweet! I'm great with numbers!"
else
	puts "That's a cool major!"
end


puts "What is your major?"
major = gets.chomp

# Set default response
major_responses = Hash.new("That's a cool major!")

# Add other responses
major_responses["Biology"] = "Mmm the study of life itself!"
major_responses["Computer Science"] = "I'm a computer!"
major_responses["English"] = "No way! What's your favorite book?"
major_responses["Math"] = "Sweet! I'm great with numbers!"

puts major_responses[major]